module Test.Main where

import Prelude
import Effect (Effect)
import Effect.Console (log)
import Foreign (Foreign, typeOf, tagOf, isNull, isUndefined, isArray, unsafeToForeign, unsafeFromForeign, readString, readBoolean, readNumber, readInt, readArray)
import Test.Assert (assert)
import Control.Monad.Except (runExcept)
import Data.Either (Either(..), isRight, isLeft)
import Foreign.Index (readIndex, readProp)

-- These cross the ordinary FFI as opaque native Foreign values in Go.
foreign import nativeValue :: String -> Foreign

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
  assert $ runExcept (readNumber (unsafeToForeign 42)) == Right 42.0
  assert $ runExcept (readNumber (unsafeToForeign (-42))) == Right (-42.0)
  assert $ isLeft $ runExcept $ readNumber (unsafeToForeign "42")

  log "Testing readInt"
  assert $ isRight $ runExcept $ readInt (unsafeToForeign 42)
  assert $ isRight $ runExcept $ readInt (unsafeToForeign 42.0)
  assert $ runExcept (readInt (unsafeToForeign 42)) == Right 42
  assert $ isLeft $ runExcept $ readInt (unsafeToForeign 42.5)
  assert $ isLeft $ runExcept $ readInt (unsafeToForeign "42")

  log "Testing readArray"
  assert $ isRight $ runExcept $ readArray (unsafeToForeign [1, 2, 3])
  assert $ isLeft $ runExcept $ readArray (unsafeToForeign 42)

  log "Testing native Foreign readers"
  assert $ runExcept (readString (nativeValue "string")) == Right "hello"
  assert $ runExcept (readString (nativeValue "emptyString")) == Right ""
  assert $ runExcept (readString (nativeValue "namedString")) == Right "named"
  assert $ runExcept (readBoolean (nativeValue "true")) == Right true
  assert $ runExcept (readBoolean (nativeValue "false")) == Right false
  assert $ runExcept (readBoolean (nativeValue "namedBoolean")) == Right true
  assert $ runExcept (readNumber (nativeValue "namedNumber")) == Right 42.5
  assert $ runExcept (readNumber (nativeValue "namedInteger")) == Right 42.0
  assert $ runExcept (readNumber (nativeValue "unsigned")) == Right 42.0
  assert $ runExcept (readNumber (nativeValue "smallSigned")) == Right (-42.0)
  assert $ runExcept (readNumber (nativeValue "float32")) == Right 42.5
  assert $ isLeft $ runExcept $ readString (nativeValue "true")
  assert $ isLeft $ runExcept $ readBoolean (nativeValue "string")
  assert $ isLeft $ runExcept $ readArray (nativeValue "object")
  assert $ isLeft $ runExcept $ readString (nativeValue "null")
  assert $ isLeft $ runExcept $ readBoolean (nativeValue "null")
  assert $ isLeft $ runExcept $ readArray (nativeValue "null")
  assert $ case runExcept (readArray (nativeValue "array")) of
    Right [ text, flag, number ] ->
      runExcept (readString text) == Right "hello" &&
      runExcept (readBoolean flag) == Right true &&
      runExcept (readInt number) == Right 42
    _ -> false
  assert $ case runExcept (readArray (nativeValue "typedArray")) of
    Right [ first, second ] ->
      runExcept (readString first) == Right "a" &&
      runExcept (readString second) == Right "b"
    _ -> false
  assert $ case runExcept (readArray (nativeValue "namedArray")) of
    Right [ first, second ] ->
      runExcept (readInt first) == Right 1 &&
      runExcept (readInt second) == Right 2
    _ -> false
  assert $ case runExcept (readArray (nativeValue "nested")) of
    Right [ inner ] -> case runExcept (readArray inner) of
      Right [ text, flag ] ->
        runExcept (readString text) == Right "nested" &&
        runExcept (readBoolean flag) == Right false
      _ -> false
    _ -> false
  assert $ case runExcept (readArray (nativeValue "emptyArray")) of
    Right [] -> true
    _ -> false
  assert $ case runExcept (readArray (nativeValue "nilArray")) of
    Right [] -> true
    _ -> false
  assert $ case runExcept (readArray (unsafeToForeign [ unsafeToForeign (\value -> value + 1) ])) of
    Right [ value ] -> (unsafeFromForeign value :: Int -> Int) 41 == 42
    _ -> false

  log "Testing native Foreign indices"
  assert $ runExcept (readIndex 0 (nativeValue "array") >>= readString) == Right "hello"
  assert $ runExcept (readIndex 1 (nativeValue "array") >>= readBoolean) == Right true
  assert $ runExcept (readIndex 2 (nativeValue "array") >>= readInt) == Right 42
  assert $ runExcept (readIndex 1 (nativeValue "typedArray") >>= readString) == Right "b"
  assert $ runExcept (readIndex 1 (nativeValue "namedArray") >>= readInt) == Right 2
  assert $ runExcept (readIndex 0 (nativeValue "nested") >>= readIndex 0 >>= readString) == Right "nested"
  assert $ runExcept (readProp "1" (nativeValue "typedArray") >>= readString) == Right "b"
  assert $ runExcept (readProp "value" (nativeValue "object") >>= readInt) == Right 42
  assert $ runExcept (isUndefined <$> readIndex (-1) (nativeValue "array")) == Right true
  assert $ runExcept (isUndefined <$> readIndex 3 (nativeValue "array")) == Right true
  assert $ runExcept (isUndefined <$> readIndex 0 (nativeValue "nilArray")) == Right true
  assert $ runExcept (isUndefined <$> readProp "not-an-index" (nativeValue "array")) == Right true
  assert $ runExcept (isUndefined <$> readProp "01" (nativeValue "array")) == Right true

  log "All tests passed"
