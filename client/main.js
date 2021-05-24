import { generator } from "./modules/generator.js";
import { generatorTest } from "../test/generatorTest.js";

let mm = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30", "31", "32", "33", "34", "35", "36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46", "47", "48", "49", "50", "51", "52", "53", "54", "55", "56", "57", "58", "59", "60", "61", "62", "63", "64", "65", "66", "67", "68", "69", "70", "71", "72", "73", "74", "75", "76", "77", "78", "79", "80", "81", "82", "83", "84", "85", "86", "87", "88", "89", "90", "91", "92", "93", "94", "95", "96", "97", "98", "99", "100", "101", "102", "103", "104", "105", "106", "107", "108", "109"]
generatorTest();
var circleWidth = 5
const palette = {
  lightgray: "#D9DEDE",
  gray: "#C3C8C8",
  mediumgray: "#536870",
  orange: "#BD3613",
  purple: "#595AB7",
  yellowgreen: "#738A05",
};
var charge = -45
let forceWidth = $("#graphic svg").width();
let forceHeight = $("#graphic svg").height();
let aWidth = $(".overviewpanel").width();
let aHeight = $(".overviewpanel").height();

function buildBarChart(data) {
  let delay = 250;
  let yearStep = 1;
  let yearMin = d3.min(data, (d) => d.year);//1841 
  let yearMax = d3.max(data, (d) => d.year);//2019
  let width = aWidth;
  let height = aHeight;
  console.log(aWidth, aHeight)
  let margin = { top: 20, right: 40, bottom: 30, left: 20 };
  // console.log('123', Array.from(d3.group(data, (d) => d.age).keys()).sort(d3.ascending))
  let x = d3
    .scaleBand()
    .domain(mm)
    .range([width - margin.right, margin.left]);//[w - 30, 0]

  // console.log(x(1),x(10),x(20),x(30),x(40),x(50))

  // console.log(Array.from(d3.group(data, (d) => d.age).keys()).sort(d3.ascending))
  console.log('max value', d3.ticks(...d3.extent(data, (d) => d.age), width / 40))
  let y = d3
    .scaleLinear()
    .domain([0, 3270])
    .range([height - margin.bottom, margin.top]);

  let color = d3.scaleOrdinal(["M", "F"], ["#4e79a7", "#e15759"]);
  let xAxis = (g) =>
    g
      .attr("transform", `translate(0,${height - margin.bottom})`)
      .call(
        d3
          .axisBottom(x)
          // .tickValues(d3.ticks(...d3.extent(data, (d) => d.age), width / 40))
          .tickValues([0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95])
          .tickSizeOuter(0)//
      )
      .call((g) =>
        g
          .append("text")
          .attr("x", 45)
          .attr("y", margin.bottom - 8)
          .attr("fill", "currentColor")
          .attr("text-anchor", "end")
          // .text("← Age")
          .text("← Mins")
      );
  let yAxis = (g) =>
    g
      .attr("transform", `translate(${width - margin.right},0)`)
      .call(d3.axisRight(y).ticks(null, "s"))
      .call((g) => g.select(".domain").remove())
      .call((g) =>
        g
          .append("text")
          .attr("x", margin.right)
          .attr("y", 10)
          .attr("fill", "currentColor")
          .attr("text-anchor", "end")
          // .text("Population ↑")
          .text("TotalIn ↑")
      );
  let svg = d3
    .select(".overviewpanel svg")
    .attr("viewBox", [0, 0, width, height]);
  // aWidth aHeight
  svg.append("g").call(xAxis);
  svg.append("g").call(yAxis);

  let group = svg.append("g");
  let rect = group.selectAll("rect");
  let group2 = svg.append("g")
  let otherrect = group2.selectAll("rect");

  function update() {
    let year = +$("#year-slider").val();//选择的年份

    // debugger;    x.step()  该函数返回相邻频段起点之间的距离。
    const dx = (x.step() * (year - yearMin)) / yearStep;

    const t = svg.transition().ease(d3.easeLinear).duration(100);
    // debugger;
    rect = rect
      .data(
        data.filter((d) => +d.year === year),
        (d) => `${d.sex}:${d.year - d.age}`
      )//                           ???
      .join(
        (enter) =>
          // console.log('123')
          enter
            .append("rect")
            .style("mix-blend-mode", "darken")
            .attr("fill", (d) => color('F'))
            .attr("x", (d) => x(d.age) + dx)
            .attr("y", (d) => y(0))
            .attr("width", x.bandwidth() + 1)
            .attr("height", 0),
        (update) => update,
        (exit) =>
          exit.call((rect) =>
            rect.transition(t).remove().attr("y", y(0)).attr("height", 0)
          )
      );
    //", "#e15759
    otherrect = otherrect
      .data(
        data.filter((d) => +d.year === year),
        (d) => `${d.sex}:${d.year - d.age}`
      )//                           ???
      .join(
        (enter) =>
          // console.log('123')
          enter
            .append("rect")
            .style("mix-blend-mode", "darken")
            .attr("fill", (d) => color('M'))
            .attr("x", (d) => {
              // console.log('123',d)
              return x(d.age) + dx + x.step()
            })
            .attr("y", (d) => y(0))
            .attr("width", x.bandwidth() + 1)
            .attr("height", 0)
        ,
        (update) => update,
        (exit) =>
          exit.call((rect) =>
            rect.transition(t).remove().attr("y", y(0)).attr("height", 0)
          )
      );
    otherrect
      .transition(t)
      .attr("y", (d) => y(d.value))
      .attr("height", (d) => {
        if(d.age === '0'){
          console.log('123',d)
          return 0
        }
        return y(0) - y(d.value)
      });
    rect
      .append("title")
      // .text(d => {
      //   return ""d.
      //   ${d.sex}${year - d.year + d.age}+d.value
      // })
      .html((d) => {
        return `<span>time:-${d.age}</span><span>Value:${d.value}</span>`
      });
    rect
      .transition(t)
      .attr("y", (d) => y(d.value))
      .attr("height", (d) => y(0) - y(d.value));

    group2.transition(t).attr("transform", `translate(${-dx},0)`);

    group.transition(t).attr("transform", `translate(${-dx},0)`);
  }
  $("#year-slider").attr("min", yearMin);
  $("#year-slider").attr("max", yearMax);
  $("#year-slider").val(yearMax);
  $("#year-slider").change(update);
}

d3.csv("./assets/icelandic-population.csv").then((data) => {
  // console.log(data);
  buildBarChart(data);
});



let drag = simulation => {

  function dragstarted(event) {
    if (!event.active) simulation.alphaTarget(0.3).restart();
    event.subject.fx = event.subject.x;
    event.subject.fy = event.subject.y;
  }

  function dragged(event) {
    event.subject.fx = event.x;
    event.subject.fy = event.y;
  }

  function dragended(event) {
    if (!event.active) simulation.alphaTarget(0);
    event.subject.fx = null;
    event.subject.fy = null;
  }

  return d3.drag()
    .on("start", dragstarted)
    .on("drag", dragged)
    .on("end", dragended);
}

function buildforce() {
  // Generate test data， 每个节点随机连接到x个节点，x在0~Maxcon中随机。
  let nodes = [];
  let numNodes = 100;
  let Maxcon = 10 //每个节点的
  for (let x = 0; x < numNodes; x++) {
    let targetAry = [];
    let connections = (Math.round(Math.random() * Maxcon));
    for (let y = 0; y < connections; y++) {
      targetAry.push(Math.floor(Math.random() * numNodes))
    }
    nodes.push({
      id: x,
      name: "node_" + x,
      target: targetAry
    })
  }

  // Create the links array from the generated data
  var links = [];
  for (var i = 0; i < nodes.length; i++) {
    if (nodes[i].target !== undefined) {
      for (var j = 0; j < nodes[i].target.length; j++) {
        links.push({
          source: i,
          target: nodes[i].target[j],
          sourcenode: nodes[i],
          targetnode: nodes[nodes[i].target[j]]
        })
      }
    }
  }
  // links nodes
  const simulation = d3.forceSimulation(nodes)
    .force("charge", d3.forceManyBody().strength(-45))
    .force("link", d3.forceLink(links))
    .force("center", d3.forceCenter(forceWidth / 2, forceHeight / 2));

  var fdGraph = d3.select('#graphic svg')
    .attr('viewBox', [0, 0, forceWidth, forceHeight])


  const link = fdGraph.append("g")
    .attr('stroke', palette.gray)
    .attr("stroke-opacity", 0.6)
    .selectAll("line")
    .data(links)
    .join("line")
    .attr("stroke-width", d => Math.sqrt(d.value))
    .attr('class', function (d, i) {
      // Add classes to lines to identify their from's and to's
      var theClass = 'from_' + d.sourcenode.id + ' line ';
      if (d.target !== undefined) {
        theClass += 'to_' + d.targetnode.id
      }
      // line_1 line to_2 (ps:这里表示分别属于三个类，line1、line、to_2)
      return theClass
    })
  // .on('mouseover', console.log)

  // Create the SVG circles for the nodes
  function overnode(event, node) {
    //重置
    d3.selectAll('.node')
      .attr('r', 5)
    d3.selectAll('.line')
      .attr('stroke', palette.lightgray)
      .attr('stroke-width', 1)

    // Hightlight the nodes that the current node connects to
    // for (let i = 0; i < node.target.length; i++) {
    //   d3.select('#node_' + d.target[i]).select('text')
    //     .attr('font-size', '14')
    //     .attr('font-weight', 'bold')
    // }
    let nodeid = node.id
    // Node
    // link connect to this node is orange
    // link connect from this node to other is purple
    for (let x = 0; x < links.length; x++) {
      // 高亮我连出去的node
      if (links[x].sourcenode.id === nodeid) {
        d3.select('#node_' + links[x].targetnode.id)
          .attr('r', 10)
      }
      if (links[x].targetnode.id === nodeid) {
        d3.select('#node_' + links[x].sourcenode.id)
          .attr('r', 10)
      }
    }
    //高亮我
    d3.select("#node_" + nodeid)
      .attr('r', 10)

    // Link
    // Highlight the connections to this node高亮到我的边
    d3.selectAll('.to_' + nodeid)
      .attr('stroke', palette.orange)
      .attr('stroke-width', 3)

    // Highlight the connections from this node高亮从我出去的边
    d3.selectAll('.from_' + nodeid)
      .attr('stroke', palette.purple)
      .attr('stroke-width', 3)

    // When mousing over a node, 
    // make it more repulsive so it stands out 让他排斥周围的节点而突出
    console.log(node)
    function strength(node2, id) {
      // console.log("node",node,"id",id)
      if (id != nodeid) {
        // Make the nodes connected to more repulsive
        for (let i = 0; i < node.target.length; i++) {
          if (id === node.target[i]) {
            return charge * 8
          }
        }
        // Make the nodes connected from more repulsive
        for (let x = 0; x < links.length; x++) {
          if (links[x].sourcenode.id === id) {
            if (links[x].targetnode.id === nodeid) {
              return charge * 8
            }
          }
        }
        // 不相关的节点保持原来的样子
        return charge * 1;
      } else {
        // Make the selected node more repulsive
        return charge * 10;
      }
    }
    simulation.force("charge", d3.forceManyBody().strength(strength));
    simulation.restart();
  }
  const node = fdGraph.append("g")
    .attr("stroke", "#fff")
    .attr("stroke-width", 1.5)
    .selectAll("circle")
    .data(nodes)
    .join("circle")
    .attr('r', circleWidth)
    .attr('class', 'node')
    .attr('id', function (d) { return "node_" + d.id })
    .attr('fill', function (d, i) {
      // Color 1/3 of the nodes each color
      // Depending on the data, this can be made more meaningful
      // 可以根据节点的压力来改变颜色。
      if (i < (numNodes / 3)) {
        return palette.orange
      } else if (i < (numNodes - (numNodes / 3))) {
        return palette.purple
      }
      return palette.yellowgreen
    })
    .on('mouseover', function overnode(event, node) {
      //重置
      d3.selectAll('.node')
        .attr('r', 5)
      d3.selectAll('.line')
        .attr('stroke', palette.lightgray)
        .attr('stroke-width', 1)

      // Hightlight the nodes that the current node connects to
      // for (let i = 0; i < node.target.length; i++) {
      //   d3.select('#node_' + d.target[i]).select('text')
      //     .attr('font-size', '14')
      //     .attr('font-weight', 'bold')
      // }
      let nodeid = node.id
      // Node
      // link connect to this node is orange
      // link connect from this node to other is purple
      for (let x = 0; x < links.length; x++) {
        // 高亮我连出去的node
        if (links[x].sourcenode.id === nodeid) {
          d3.select('#node_' + links[x].targetnode.id)
            .attr('r', 10)
        }
        if (links[x].targetnode.id === nodeid) {
          d3.select('#node_' + links[x].sourcenode.id)
            .attr('r', 10)
        }
      }
      //高亮我
      d3.select("#node_" + nodeid)
        .attr('r', 10)

      // Link
      // Highlight the connections to this node高亮到我的边
      d3.selectAll('.to_' + nodeid)
        .attr('stroke', palette.orange)
        .attr('stroke-width', 3)

      // Highlight the connections from this node高亮从我出去的边
      d3.selectAll('.from_' + nodeid)
        .attr('stroke', palette.purple)
        .attr('stroke-width', 3)

      // When mousing over a node, 
      // make it more repulsive so it stands out 让他排斥周围的节点而突出
      console.log(node)
      function strength(node2, id) {
        // console.log("node",node,"id",id)
        if (id != nodeid) {
          // Make the nodes connected to more repulsive
          for (let i = 0; i < node.target.length; i++) {
            if (id === node.target[i]) {
              return charge * 8
            }
          }
          // Make the nodes connected from more repulsive
          for (let x = 0; x < links.length; x++) {
            if (links[x].sourcenode.id === id) {
              if (links[x].targetnode.id === nodeid) {
                return charge * 8
              }
            }
          }
          // 不相关的节点保持原来的样子
          return charge * 1;
        } else {
          // Make the selected node more repulsive
          return charge * 10;
        }
      }
      simulation.force("charge", d3.forceManyBody().strength(strength));
      simulation.restart();
    })
    .on('click', function (event, node) { })
    .call(drag(simulation));

  //鼠标悬停时候显示
  node.append("title")
    .text(d => d.id);

  node.append('text')
    .text(function (d) {
      return d.name
    })
    .attr('font-size', '12')

  simulation.on("tick", () => {
    // Move the nodes base on charge and gravity
    link
      .attr("x1", d => d.source.x)
      .attr("y1", d => d.source.y)
      .attr("x2", d => d.target.x)
      .attr("y2", d => d.target.y);

    node
      .attr("cx", d => d.x)
      .attr("cy", d => d.y);
  });

  // 聚焦在34号的方法如下
  // overnode(null, nodes[34])

}

buildforce()