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