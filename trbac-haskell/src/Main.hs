{-# LANGUAGE OverloadedStrings #-}
{-# LANGUAGE RankNTypes #-}
{-# LANGUAGE ScopedTypeVariables #-}
{-# LANGUAGE ImpredicativeTypes #-}

module Main where

import qualified Data.Map as Map
import qualified Data.Set as Set
import Data.Text (Text, pack)
import TRBAC.Types
import TRBAC.Config
import TRBAC.Context
import TRBAC.ConstraintRunner
import TRBAC.Authorization

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
      
      let constraintMap :: Map.Map Constraint (forall c. Context c => c -> Bool)
          constraintMap = Map.fromList
            [ (businessHoursConstraint, \_ -> True) -- Always pass for this example
            , (auditLogConstraint, \_ -> True) -- Simplified for type compatibility
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
      
      -- Log the access attempt
      putStrLn $ "AUDIT: " ++ show (getAction ctx) ++ " " ++
        show (getResourceType ctx) ++ " by " ++ show (getRoles ctx)
      
      if authorized
        then putStrLn "Access granted"
        else putStrLn "Access denied"