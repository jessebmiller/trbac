module TRBAC.Context where

import TRBAC.Types

-- | Type class for authorization contexts
class Context c where
  -- | Get the action being performed
  getAction :: c -> Action
  
  -- | Get the type of resource being accessed
  getResourceType :: c -> ResourceType
  
  -- | Get all roles assigned to the actor
  getRoles :: c -> [Role]

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