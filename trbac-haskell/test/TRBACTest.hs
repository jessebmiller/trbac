module Main (main) where

import Test.Hspec
import Test.QuickCheck
import qualified Data.Map as Map
import qualified Data.Set as Set

import TRBAC.Types
import TRBAC.Context
import TRBAC.ConstraintRunner
import TRBAC.Authorization

-- | Generate arbitrary actions
instance Arbitrary Action where
  arbitrary = Action <$> elements ["read", "write", "create", "delete"]

-- | Generate arbitrary resource types
instance Arbitrary ResourceType where
  arbitrary = ResourceType <$> elements ["document", "folder", "user"]

-- | Generate arbitrary roles
instance Arbitrary Role where
  arbitrary = Role <$> elements ["admin", "writer", "reader"]

-- | Generate arbitrary permissions
instance Arbitrary Permission where
  arbitrary = do
    actions <- listOf1 arbitrary
    resourceTypes <- listOf1 arbitrary
    constraints <- listOf arbitrary
    return Permission
      { permActions = Set.fromList actions
      , permResourceTypes = Set.fromList resourceTypes
      , permConstraints = Set.fromList constraints
      }

main :: IO ()
main = hspec $ do
  describe "TRBAC.may" $ do
    it "grants access when a relevant permission exists with no constraints" $ do
      property $ \action resourceType role -> do
        let perm = Permission
              { permActions = Set.singleton action
              , permResourceTypes = Set.singleton resourceType
              , permConstraints = Set.empty
              }
        let privileges = Map.singleton role [perm]
        let ctx = BasicContext
              { basicAction = action
              , basicResourceType = resourceType
              , basicRoles = [role]
              }
        let runner = FunctionMapConstraintRunner Map.empty
        
        may privileges runner ctx `shouldBe` True
    
    it "denies access when no relevant permission exists" $ do
      property $ \action resourceType role -> do
        let perm = Permission
              { permActions = Set.singleton action
              , permResourceTypes = Set.singleton resourceType
              , permConstraints = Set.empty
              }
        let privileges = Map.singleton role [perm]
        let ctx = BasicContext
              { basicAction = Action "different"
              , basicResourceType = resourceType
              , basicRoles = [role]
              }
        let runner = FunctionMapConstraintRunner Map.empty
        
        may privileges runner ctx `shouldBe` False
    
    it "denies access when constraints fail" $ do
      property $ \action resourceType role -> do
        let constraint = Constraint "always_fail"
        let perm = Permission
              { permActions = Set.singleton action
              , permResourceTypes = Set.singleton resourceType
              , permConstraints = Set.singleton constraint
              }
        let privileges = Map.singleton role [perm]
        let ctx = BasicContext
              { basicAction = action
              , basicResourceType = resourceType
              , basicRoles = [role]
              }
        let runner = FunctionMapConstraintRunner $ Map.singleton constraint (const False)
        
        may privileges runner ctx `shouldBe` False

    -- Additional test cases for comprehensive coverage
    it "denies access when roles are empty" $ do
      property $ \action resourceType -> do
        let perm = Permission
              { permActions = Set.singleton action
              , permResourceTypes = Set.singleton resourceType
              , permConstraints = Set.empty
              }
        let privileges = Map.empty
        let ctx = BasicContext
              { basicAction = action
              , basicResourceType = resourceType
              , basicRoles = []
              }
        let runner = FunctionMapConstraintRunner Map.empty
        
        may privileges runner ctx `shouldBe` False

    it "denies access when permissions are empty" $ do
      property $ \action resourceType role -> do
        let privileges = Map.singleton role []
        let ctx = BasicContext
              { basicAction = action
              , basicResourceType = resourceType
              , basicRoles = [role]
              }
        let runner = FunctionMapConstraintRunner Map.empty
        
        may privileges runner ctx `shouldBe` False

    it "handles multiple roles with overlapping permissions" $ do
      property $ \action resourceType role1 role2 -> do
        let perm1 = Permission
              { permActions = Set.singleton action
              , permResourceTypes = Set.singleton resourceType
              , permConstraints = Set.empty
              }
        let perm2 = Permission
              { permActions = Set.singleton (Action "different")
              , permResourceTypes = Set.singleton resourceType
              , permConstraints = Set.empty
              }
        let privileges = Map.fromList [(role1, [perm1]), (role2, [perm2])]
        let ctx = BasicContext
              { basicAction = action
              , basicResourceType = resourceType
              , basicRoles = [role1, role2]
              }
        let runner = FunctionMapConstraintRunner Map.empty
        
        may privileges runner ctx `shouldBe` True

    -- New tests for complex constraints
    it "handles complex constraints correctly" $ do
      property $ \action resourceType role -> do
        let complexConstraint = Constraint "complex"
        let perm = Permission
              { permActions = Set.singleton action
              , permResourceTypes = Set.singleton resourceType
              , permConstraints = Set.singleton complexConstraint
              }
        let privileges = Map.singleton role [perm]
        let ctx = BasicContext
              { basicAction = action
              , basicResourceType = resourceType
              , basicRoles = [role]
              }
        let runner = FunctionMapConstraintRunner $ Map.singleton complexConstraint (\_ -> True)
        
        may privileges runner ctx `shouldBe` True

    it "handles invalid inputs gracefully" $ do
      let invalidAction = Action ""
      let invalidResourceType = ResourceType ""
      let invalidRole = Role ""
      let perm = Permission
            { permActions = Set.singleton invalidAction
            , permResourceTypes = Set.singleton invalidResourceType
            , permConstraints = Set.empty
            }
      let privileges = Map.singleton invalidRole [perm]
      let ctx = BasicContext
            { basicAction = invalidAction
            , basicResourceType = invalidResourceType
            , basicRoles = [invalidRole]
            }
      let runner = FunctionMapConstraintRunner Map.empty
      
      may privileges runner ctx `shouldBe` False

    it "performs well with large numbers of roles and permissions" $ do
      property $ \action resourceType -> do
        let roles = map (Role . pack . show) [1..1000]
        let perms = replicate 1000 Permission
              { permActions = Set.singleton action
              , permResourceTypes = Set.singleton resourceType
              , permConstraints = Set.empty
              }
        let privileges = Map.fromList $ zip roles (repeat perms)
        let ctx = BasicContext
              { basicAction = action
              , basicResourceType = resourceType
              , basicRoles = roles
              }
        let runner = FunctionMapConstraintRunner Map.empty
        
        may privileges runner ctx `shouldBe` True

    it "logs errors and handles exceptions" $ do
      let faultyConstraint = Constraint "faulty"
      let perm = Permission
            { permActions = Set.singleton (Action "read")
            , permResourceTypes = Set.singleton (ResourceType "document")
            , permConstraints = Set.singleton faultyConstraint
            }
      let privileges = Map.singleton (Role "reader") [perm]
      let ctx = BasicContext
            { basicAction = Action "read"
            , basicResourceType = ResourceType "document"
            , basicRoles = [Role "reader"]
            }
      let runner = FunctionMapConstraintRunner $ Map.singleton faultyConstraint (\_ -> error "Constraint evaluation failed")
      
      evaluate (may privileges runner ctx) `shouldThrow` anyErrorCall
