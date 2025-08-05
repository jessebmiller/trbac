{-# LANGUAGE RankNTypes #-}
{-# LANGUAGE ImpredicativeTypes #-}

module TRBAC.ConstraintRunner where

import qualified Data.Map as Map
import TRBAC.Types
import TRBAC.Context

-- | Type class for constraint runners
class ConstraintRunner r where
  -- | Evaluate a constraint in a context
  runConstraint :: Context c => r -> Constraint -> c -> Bool

-- | A constraint runner using a map of functions
newtype FunctionMapConstraintRunner = FunctionMapConstraintRunner
  { getFunctionMap :: Map.Map Constraint (forall c. Context c => c -> Bool)
  }

instance ConstraintRunner FunctionMapConstraintRunner where
  runConstraint runner constraint ctx =
    case Map.lookup constraint (getFunctionMap runner) of
      Just f -> f ctx
      Nothing -> False
