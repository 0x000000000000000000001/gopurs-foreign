export function nativeValue(name) {
  switch (name) {
    case "string": return "hello";
    case "emptyString": return "";
    case "namedString": return "named";
    case "true": return true;
    case "false": return false;
    case "namedBoolean": return true;
    case "namedNumber": return 42.5;
    case "namedInteger": return 42;
    case "unsigned": return 42;
    case "smallSigned": return -42;
    case "float32": return 42.5;
    case "array": return ["hello", true, 42];
    case "typedArray": return ["a", "b"];
    case "namedArray": return [1, 2];
    case "nested": return [["nested", false]];
    case "emptyArray": return [];
    case "nilArray": return [];
    case "object": return { value: 42 };
    case "null": return null;
    default: throw new Error("Unknown native Foreign test value: " + name);
  }
}
