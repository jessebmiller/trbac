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

instance FromJSONKey Role where
  fromJSONKey = FromJSONKeyText Role

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