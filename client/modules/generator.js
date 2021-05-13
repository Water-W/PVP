export function randomString(l, r) {
  let chars = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678";
  let length = Math.floor(Math.random() * (r - l)) + l;
  let str = "";
  for (let i = 0; i < length; ++i) {
    let pos = Math.floor(Math.random() * chars.length);
    str = str + chars.charAt(pos);
  }
  return str;
}

export function randomName(num, l, r) {
  let nameDict = {};
  let nameList = [];
  for (let i = 0; i < num; ++i) {
    let name = randomString(l, r);
    if (!(name in nameDict)) {
      nameDict[name] = true;
      nameList.push(name);
    } else {
      i--;
    }
  }
  return nameList;
}

export function generator(nodeNum, protocolNum) {
  let nodeName = randomName(nodeNum, 5, 5);
  let linkList = [];
  let linkNum = 0;
  for (let i of nodeName) {
    for (let j of nodeName) {
      if (Math.random() < 1 / 5) {
        linkList.push({ id: linkNum, source: i, target: j });
        linkNum += 1;
      }
    }
  }
  function stepGenerator() {}
  return { node: nodeName, link: linkList, generator: stepGenerator };
}
