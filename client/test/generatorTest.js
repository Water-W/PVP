import { randomString, randomName, generator } from "../modules/generator.js";

export function generatorTest() {
  let strs = randomName(50, 3, 5);
  for (let i of strs) {
    if (i.length < 3 || i.length > 5) {
      console.log("test error - ", i);
    }
  }
}
