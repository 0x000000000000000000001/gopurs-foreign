module Test.Main where

import Prelude
import Effect (Effect)
import Effect.Console (log)
import Foreign (typeOf, tagOf, isNull, isUndefined, isArray, unsafeToForeign, readString, readBoolean, readNumber, readInt, readArray)
import Test.Assert (assert)
import Control.Monad.Except (runExcept)
import Data.Either (isRight, isLeft)

main :: Effect Unit
main = do
  log "Testing typeOf"
  assert $ typeOf (unsafeToForeign "hello") == "string"
  assert $ typeOf (unsafeToForeign 42) == "number"
  assert $ typeOf (unsafeToForeign 42.5) == "number"
  assert $ typeOf (unsafeToForeign true) == "boolean"
  assert $ typeOf (unsafeToForeign (\x -> x)) == "function"
  assert $ typeOf (unsafeToForeign { a: 1 }) == "object"
  assert $ typeOf (unsafeToForeign [ 1, 2, 3 ]) == "object"

  log "Testing tagOf"
  assert $ tagOf (unsafeToForeign "hello") == "String"
  assert $ tagOf (unsafeToForeign 42) == "Number"
  assert $ tagOf (unsafeToForeign 42.5) == "Number"
  assert $ tagOf (unsafeToForeign true) == "Boolean"
  assert $ tagOf (unsafeToForeign (\x -> x)) == "Function"
  assert $ tagOf (unsafeToForeign { a: 1 }) == "Object"
  assert $ tagOf (unsafeToForeign [ 1, 2, 3 ]) == "Array"

  log "Testing isArray"
  assert $ isArray (unsafeToForeign [ 1, 2, 3 ]) == true
  assert $ isArray (unsafeToForeign 42) == false
  assert $ isArray (unsafeToForeign { a: 1 }) == false
  assert $ isArray (unsafeToForeign "test") == false

  log "Testing isNull and isUndefined"
  assert $ isNull (unsafeToForeign "hello") == false
  assert $ isNull (unsafeToForeign { a: 1 }) == false
  assert $ isUndefined (unsafeToForeign "hello") == false
  assert $ isUndefined (unsafeToForeign { a: 1 }) == false

  log "Testing readString"
  assert $ isRight $ runExcept $ readString (unsafeToForeign "hello")
  assert $ isLeft $ runExcept $ readString (unsafeToForeign 42)

  log "Testing readBoolean"
  assert $ isRight $ runExcept $ readBoolean (unsafeToForeign true)
  assert $ isLeft $ runExcept $ readBoolean (unsafeToForeign 42)

  log "Testing readNumber"
  assert $ isRight $ runExcept $ readNumber (unsafeToForeign 42.5)
  assert $ isRight $ runExcept $ readNumber (unsafeToForeign 42)
  assert $ isLeft $ runExcept $ readNumber (unsafeToForeign "42")

  log "Testing readInt"
  assert $ isRight $ runExcept $ readInt (unsafeToForeign 42)
  assert $ isRight $ runExcept $ readInt (unsafeToForeign 42.0)
  assert $ isLeft $ runExcept $ readInt (unsafeToForeign 42.5)
  assert $ isLeft $ runExcept $ readInt (unsafeToForeign "42")

  log "Testing readArray"
  assert $ isRight $ runExcept $ readArray (unsafeToForeign [1, 2, 3])
  assert $ isLeft $ runExcept $ readArray (unsafeToForeign 42)

  log "All tests passed"
