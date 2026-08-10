module Scratch where

import Prelude
import Foreign (unsafeToForeign, readString, readNumber, readInt, readBoolean, readArray)
import Control.Monad.Except (runExcept)
import Data.Either (isRight, isLeft)
import Effect (Effect)
import Effect.Console (log)

main :: Effect Unit
main = do
  log $ show $ isRight $ runExcept $ readString (unsafeToForeign "hello")
  log $ show $ isRight $ runExcept $ readString (unsafeToForeign 42)
  log $ show $ isRight $ runExcept $ readInt (unsafeToForeign 42)
