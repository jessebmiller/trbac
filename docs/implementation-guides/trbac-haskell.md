# TRBAC Implementation Guide for Haskell

This guide provides a practical approach to implementing TRBAC (Typed Role-Based Access Control) in Haskell. It follows the formal specification while leveraging Haskell's strengths: strong static typing, pure functions, and expressive type system.

## Type Definitions

The Haskell implementation of TRBAC utilizes algebraic data types and type classes to model the core concepts:

### Core Types

```haskell
module TRBAC.Types where

import qualified Data.Set as Set
import qualified Data.Map as Map
import Data.Text (Text)

-- | Represents an action that can be taken on a resource
newtype Action = Action Text
  deriving (Eq, Ord, Show)

-- | Represents a type of resource that can be protected
newtype ResourceType = ResourceType Text
  deriving (Eq, Ord, Show)

-- | Represents a role that can be assigned to an actor
newtype Role = Role Text
  deriving (Eq, Ord, Show)

-- | Represents a constraint that can restrict permissions
newtype Constraint = Constraint Text
  deriving (Eq, Ord, Show)

-- | Represents a permission to take actions on resource types
data Permission = Permission
  { permActions :: Set.Set Action
  , permResourceTypes :: Set.Set ResourceType
  , permConstraints :: Set.Set Constraint
  } deriving (Eq, Show)

-- | Maps roles to their permissions
type Privileges = Map.Map Role [Permission]
```

### Context Type Class

The `Context` type class defines the interface for authorization contexts:

```haskell
-- | Type class for authorization contexts
class Context c where
  -- | Get the action being performed
  getAction :: c -> Action
  
  -- | Get the type of resource being accessed
  getResourceType :: c -> ResourceType
  
  -- | Get all roles assigned to the actor
  getRoles :: c -> [Role]
```

### Constraint Runner Type Class

The `ConstraintRunner` type class defines how constraints are evaluated:

```haskell
-- | Type class for constraint runners
class ConstraintRunner r where
  -- | Evaluate a constraint in a context
  runConstraint :: r -> Constraint -> c -> Bool
    where c :: Context c
```

## Implementation

### Basic Context Implementation

A simple implementation of the `Context` type class:

```haskell
-- | A simple context implementation
data BasicContext = BasicContext
  { basicAction :: Action
  , basicResourceType :: ResourceType
  , basicRoles :: [Role]
  }

instance Context BasicContext where
  getAction = basicAction
  getResourceType = basicResourceType
  getRoles = basicRoles
```

### Function Map Constraint Runner

A constraint runner that uses a map of functions:

```haskell
-- | A constraint runner using a map of functions
newtype FunctionMapConstraintRunner = FunctionMapConstraintRunner
  { getFunctionMap :: Map.Map Constraint (forall c. Context c => c -> Bool)
  }

instance ConstraintRunner FunctionMapConstraintRunner where
  runConstraint runner constraint ctx =
    case Map.lookup constraint (getFunctionMap runner) of
      Just f -> f ctx
      Nothing -> False
```

### Shell Command Constraint Runner

A constraint runner that executes shell commands:

```haskell
-- | A constraint runner using shell commands
data ShellCommandConstraintRunner = ShellCommandConstraintRunner
  { scriptRoot :: FilePath
  }

instance ConstraintRunner ShellCommandConstraintRunner where
  runConstraint runner constraint ctx = do
    let scriptPath = scriptRoot runner </> unpack (getConstraintName constraint)
    let args = [ unpack $ getActionName $ getAction ctx
               , unpack $ getResourceTypeName $ getResourceType ctx
               ] ++ map (unpack . getRoleName) (getRoles ctx)
    
    exitCode <- runProcess scriptPath args
    return $ exitCode == ExitSuccess
  where
    getConstraintName (Constraint name) = name
    getActionName (Action name) = name
    getResourceTypeName (ResourceType name) = name
    getRoleName (Role name) = name
```

### Authorization Function

The core authorization function implemented as a pure function:

```haskell
-- | Determines if an action is permitted in a context
may :: (Context c, ConstraintRunner r) => Privileges -> r -> c -> Bool
may privileges runner ctx =
  any permissionApplies $ relevantPermissions privileges ctx
  where
    permissionApplies perm =
      isRelevant perm ctx && allConstraintsPass perm runner ctx

-- | Checks if a permission is relevant to a context
isRelevant :: Context c => Permission -> c -> Bool
isRelevant perm ctx =
  Set.member (getAction ctx) (permActions perm) &&
  Set.member (getResourceType ctx) (permResourceTypes perm)

-- | Checks if all constraints for a permission pass
allConstraintsPass :: (Context c, ConstraintRunner r) => Permission -> r -> c -> Bool
allConstraintsPass perm runner ctx =
  all (\c -> runConstraint runner c ctx) (permConstraints perm)

-- | Gets all permissions for a set of roles
relevantPermissions :: Context c => Privileges -> c -> [Permission]
relevantPermissions privileges ctx =
  concatMap getPermsForRole (getRoles ctx)
  where
    getPermsForRole role = Map.findWithDefault [] role privileges
```

## YAML Configuration Loading

Loading privileges from YAML configuration:

```haskell
{-# LANGUAGE OverloadedStrings #-}

module TRBAC.Config where

import qualified Data.Map as Map
import qualified Data.Set as Set
import Data.Aeson
import Data.Yaml
import TRBAC.Types

instance FromJSON Action where
  parseJSON = withText "Action" $ pure . Action

instance FromJSON ResourceType where
  parseJSON = withText "ResourceType" $ pure . ResourceType

instance FromJSON Role where
  parseJSON = withText "Role" $ pure . Role

instance FromJSON Constraint where
  parseJSON = withText "Constraint" $ pure . Constraint

instance FromJSON Permission where
  parseJSON = withObject "Permission" $ \v -> Permission
    <$> (Set.fromList <$> v .: "actions")
    <*> (Set.fromList <$> v .: "resource_types")
    <*> (Set.fromList <$> v .: "constraints")

-- | Load privileges from a YAML file
loadPrivileges :: FilePath -> IO (Either ParseException Privileges)
loadPrivileges = decodeFileEither
```

## Example Usage

Here's how to use TRBAC in a Haskell application:

```haskell
module Main where

import qualified Data.Map as Map
import qualified Data.Set as Set
import Data.Text (Text, pack)
import TRBAC.Types
import TRBAC.Config

main :: IO ()
main = do
  -- Load privileges from YAML
  ePrivileges <- loadPrivileges "privileges.yaml"
  case ePrivileges of
    Left err -> putStrLn $ "Error loading privileges: " ++ show err
    Right privileges -> do
      -- Create a constraint runner
      let businessHoursConstraint = Constraint "business_hours_only"
      let auditLogConstraint = Constraint "audit_logged"
      
      let constraintMap = Map.fromList
            [ (businessHoursConstraint, const True) -- Always pass for this example
            , (auditLogConstraint, \ctx -> do
                -- Log the access attempt
                putStrLn $ "AUDIT: " ++ show (getAction ctx) ++ " " ++
                  show (getResourceType ctx) ++ " by " ++ show (getRoles ctx)
                return True
              )
            ]
      
      let runner = FunctionMapConstraintRunner constraintMap
      
      -- Create a context
      let ctx = BasicContext
            { basicAction = Action "read"
            , basicResourceType = ResourceType "document"
            , basicRoles = [Role "reader"]
            }
      
      -- Check authorization
      let authorized = may privileges runner ctx
      
      if authorized
        then putStrLn "Access granted"
        else putStrLn "Access denied"
```

## Web Framework Integration

Here's an example of how to integrate TRBAC with the Servant web framework:

```haskell
{-# LANGUAGE DataKinds #-}
{-# LANGUAGE TypeOperators #-}
{-# LANGUAGE OverloadedStrings #-}

module Main where

import Control.Monad.Reader
import Data.Maybe (fromMaybe)
import Data.Text (Text)
import qualified Data.Map as Map
import qualified Data.Set as Set
import Network.Wai
import Network.Wai.Handler.Warp
import Servant
import Servant.Server.Experimental.Auth

import TRBAC.Types
import TRBAC.Config

-- | Application state
data AppState = AppState
  { appPrivileges :: Privileges
  , appConstraintRunner :: FunctionMapConstraintRunner
  }

-- | Custom authorization type
data TRBACAuth = TRBACAuth

-- | API type definition
type API = 
  Auth '[TRBACAuth] BasicContext :> 
  "documents" :> Get '[JSON] [Text]

-- | Authorization handler
type instance AuthServerData (AuthProtect "trbac") = BasicContext

-- | Server implementation
server :: ServerT API (ReaderT AppState Handler)
server ctx = do
  -- Check authorization
  authorized <- checkAuthorization ctx
  
  if authorized
    then return ["document1", "document2"]
    else throwError err403

-- | Check authorization using TRBAC
checkAuthorization :: BasicContext -> ReaderT AppState Handler Bool
checkAuthorization ctx = do
  appState <- ask
  return $ may 
    (appPrivileges appState) 
    (appConstraintRunner appState) 
    ctx

-- | Extract context from request
contextFromRequest :: Request -> BasicContext
contextFromRequest req = BasicContext
  { basicAction = actionFromMethod (requestMethod req)
  , basicResourceType = resourceTypeFromPath (rawPathInfo req)
  , basicRoles = rolesFromRequest req
  }

-- | Convert HTTP method to Action
actionFromMethod :: Method -> Action
actionFromMethod method = Action $ case method of
  "GET" -> "read"
  "POST" -> "create"
  "PUT" -> "update"
  "DELETE" -> "delete"
  _ -> "unknown"

-- | Convert path to ResourceType
resourceTypeFromPath :: ByteString -> ResourceType
resourceTypeFromPath path 
  | path == "/documents" = ResourceType "document"
  | otherwise = ResourceType "unknown"

-- | Extract roles from request (simplified)
rolesFromRequest :: Request -> [Role]
rolesFromRequest _ = [Role "reader"] -- In practice, extract from JWT, session, etc.

-- | Main application
main :: IO ()
main = do
  -- Load privileges
  ePrivileges <- loadPrivileges "privileges.yaml"
  privileges <- case ePrivileges of
    Left err -> error $ "Failed to load privileges: " ++ show err
    Right p -> return p
  
  -- Create constraint runner
  let constraintMap = Map.fromList
        [ (Constraint "audit_logged", \_ -> return True)
        ]
  let runner = FunctionMapConstraintRunner constraintMap
  
  -- Create app state
  let appState = AppState
        { appPrivileges = privileges
        , appConstraintRunner = runner
        }
  
  -- Setup authentication
  let authHandler = mkAuthHandler $ \req -> 
        return $ contextFromRequest req
  
  let api = Proxy :: Proxy API
  let context = authHandler :. EmptyContext
  
  -- Run server
  run 8080 $ serveWithContext api context $
    hoistServerWithContext api (Proxy :: Proxy '[AuthHandler Request BasicContext]) 
      (flip runReaderT appState) server
```

## Functional Programming Advantages

Haskell's functional programming features provide several advantages for TRBAC:

### 1. Algebraic Data Types for Type Safety

Haskell's algebraic data types provide strong guarantees about the structure of data:

```haskell
-- Ensure constraints can only be created in valid ways
data Constraint 
  = BusinessHoursOnly
  | AuditLogged
  | UserOwnsResource
  | CustomConstraint Text
  deriving (Eq, Ord, Show)
```

### 2. Pattern Matching for Permission Checking

Pattern matching provides an elegant way to check permissions:

```haskell
isRelevant :: Permission -> Action -> ResourceType -> Bool
isRelevant perm action resourceType = case (action, resourceType) of
  (ReadAction, DocumentResource) -> 
    Set.member ReadAction (permActions perm) && Set.member DocumentResource (permResourceTypes perm)
  (WriteAction, DocumentResource) ->
    Set.member WriteAction (permActions perm) && Set.member DocumentResource (permResourceTypes perm)
  _ -> False
```

### 3. Monad Transformers for Context

Monad transformers can be used to elegantly compose different context sources:

```haskell
type AuthT m a = ReaderT BasicContext m a

runAuthAction :: Context c => Privileges -> ConstraintRunner r -> c -> AuthT m a -> m a
runAuthAction privileges runner ctx action = do
  let authorized = may privileges runner ctx
  if authorized
    then runReaderT action (toBasicContext ctx)
    else error "Unauthorized"
  where
    toBasicContext c = BasicContext
      { basicAction = getAction c
      , basicResourceType = getResourceType c
      , basicRoles = getRoles c
      }
```

### 4. Free Monads for Constraint DSLs

Free monads can be used to create constraint DSLs:

```haskell
-- | Constraint DSL operations
data ConstraintF next
  = CheckTime (Bool -> next)
  | CheckOwnership UserId ResourceId (Bool -> next)
  | LogAction (IO () -> next)
  deriving Functor

type ConstraintDSL = Free ConstraintF

-- | DSL Operations
checkTime :: ConstraintDSL Bool
checkTime = liftF $ CheckTime id

checkOwnership :: UserId -> ResourceId -> ConstraintDSL Bool
checkOwnership userId resourceId = liftF $ CheckOwnership userId resourceId id

logAction :: IO () -> ConstraintDSL ()
logAction action = liftF $ LogAction (const ())

-- | Example constraint using the DSL
businessHoursConstraint :: ConstraintDSL Bool
businessHoursConstraint = do
  isBusinessHours <- checkTime
  unless isBusinessHours $ logAction $ putStrLn "Access denied: outside business hours"
  return isBusinessHours
```

## Testing Strategies

Haskell's strong type system and property-based testing tools make TRBAC testing robust:

```haskell
module TRBAC.Test where

import Test.Hspec
import Test.QuickCheck
import qualified Data.Map as Map
import qualified Data.Set as Set

import TRBAC.Types

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

-- | Test authorization properties
spec :: Spec
spec = do
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
```

## Lens Integration

Haskell's lens library can be used to manipulate complex authorization structures:

```haskell
{-# LANGUAGE TemplateHaskell #-}

import Control.Lens

-- | Make lenses for Permission
makeLenses ''Permission

-- | Add an action to a permission
addAction :: Action -> Permission -> Permission
addAction action = over permActions (Set.insert action)

-- | Remove a constraint from a permission
removeConstraint :: Constraint -> Permission -> Permission
removeConstraint constraint = over permConstraints (Set.delete constraint)

-- | Update permissions for a role
updateRolePermissions :: (Permission -> Permission) -> Role -> Privileges -> Privileges
updateRolePermissions f role = over (at role . non []) (map f)
```

## Type-Safe Resource Access

Haskell's phantom types can be used to create type-safe resource access:

```haskell
{-# LANGUAGE GADTs #-}
{-# LANGUAGE DataKinds #-}
{-# LANGUAGE KindSignatures #-}

-- | Resource access mode
data AccessMode = Read | Write | Create | Delete

-- | Resource types
data ResourceTag = DocumentTag | FolderTag | UserTag

-- | Type-safe resource reference
data Resource (a :: ResourceTag) where
  Document :: Int -> Resource 'DocumentTag
  Folder :: Int -> Resource 'FolderTag
  User :: Int -> Resource 'UserTag

-- | Type-safe access request
data AccessRequest (m :: AccessMode) (a :: ResourceTag) where
  ReadRequest :: Resource a -> AccessRequest 'Read a
  WriteRequest :: Resource a -> AccessRequest 'Write a
  CreateRequest :: AccessRequest 'Create a
  DeleteRequest :: Resource a -> AccessRequest 'Delete a

-- | Convert access mode to Action
accessModeToAction :: AccessRequest m a -> Action
accessModeToAction req = case req of
  ReadRequest _ -> Action "read"
  WriteRequest _ -> Action "write"
  CreateRequest -> Action "create"
  DeleteRequest _ -> Action "delete"

-- | Convert resource tag to ResourceType
resourceTagToType :: AccessRequest m a -> ResourceType
resourceTagToType req = case req of
  ReadRequest (Document _) -> ResourceType "document"
  ReadRequest (Folder _) -> ResourceType "folder"
  ReadRequest (User _) -> ResourceType "user"
  WriteRequest (Document _) -> ResourceType "document"
  WriteRequest (Folder _) -> ResourceType "folder"
  WriteRequest (User _) -> ResourceType "user"
  CreateRequest -> ResourceType "document"  -- Example, depends on 'a'
  DeleteRequest (Document _) -> ResourceType "document"
  DeleteRequest (Folder _) -> ResourceType "folder"
  DeleteRequest (User _) -> ResourceType "user"

-- | Type-safe authorization check
authorizeRequest :: Privileges -> ConstraintRunner r -> [Role] -> AccessRequest m a -> r -> Bool
authorizeRequest privileges runner roles req r =
  let ctx = BasicContext
        { basicAction = accessModeToAction req
        , basicResourceType = resourceTagToType req
        , basicRoles = roles
        }
  in may privileges runner ctx
```

## Pure Functional Authorization Service

A pure functional approach to authorization services:

```haskell
-- | Authorization service with pure state
data AuthService = AuthService
  { authPrivileges :: Privileges
  , authConstraintRunner :: FunctionMapConstraintRunner
  }

-- | Create a new authorization service
newAuthService :: Privileges -> FunctionMapConstraintRunner -> AuthService
newAuthService = AuthService

-- | Add a permission to a role
addPermission :: Role -> Permission -> AuthService -> AuthService
addPermission role perm service =
  service { authPrivileges = Map.alter addPerm role (authPrivileges service) }
  where
    addPerm Nothing = Just [perm]
    addPerm (Just perms) = Just (perm : perms)

-- | Add a constraint function
addConstraintFunction :: Constraint -> (forall c. Context c => c -> Bool) -> AuthService -> AuthService
addConstraintFunction constraint f service =
  service { authConstraintRunner = FunctionMapConstraintRunner newMap }
  where
    oldMap = getFunctionMap (authConstraintRunner service)
    newMap = Map.insert constraint f oldMap

-- | Check authorization
checkAuthorization :: Context c => c -> AuthService -> Bool
checkAuthorization ctx service =
  may (authPrivileges service) (authConstraintRunner service) ctx
```

## Best Practices for Haskell Implementation

1. **Use the type system** to enforce correctness at compile time
2. **Leverage pure functions** for predictable authorization logic
3. **Use algebraic data types** for modeling domain concepts
4. **Define clear type class interfaces** for extensibility
5. **Use lenses** for complex nested data manipulation
6. **Implement property-based testing** to verify authorization rules
7. **Use phantom types** for type-safe resource access
8. **Apply the ReaderT pattern** for dependency injection
9. **Document all interfaces** with Haddock comments
10. **Consider performance implications** for large-scale applications
# TRBAC Implementation Guide for Haskell

This guide provides a practical approach to implementing TRBAC (Typed Role-Based Access Control) in