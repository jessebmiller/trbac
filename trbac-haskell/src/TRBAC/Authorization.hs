module TRBAC.Authorization where

import TRBAC.Types
import TRBAC.Context
import TRBAC.ConstraintRunner
import qualified Data.Set as Set
import qualified Data.Map as Map

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