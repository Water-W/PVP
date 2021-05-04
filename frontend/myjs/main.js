// Configure graphics'


infodict = {}

var width = $("#graphic").width(),
    height = $("#graphic").height();


const info = document.createElement("div");
info.setAttribute("class", "info-panel");

infopanel = d3.select(".info-panel")
    .style("position", "fixed")


console.log(width, height);

var circleWidth = 5,
    charge = -75 * 0.6,
    gravity = 0.1;

var lock_click = 0;

//给node和link编号
var nodeNum = 0
var linkNum = 0

var mylink = []
var mynode = []

var palette = {
    "lightgray": "#D9DEDE",
    "gray": "#C3C8C8",
    "mediumgray": "#536870",
    "orange": "#BD3613",
    "purple": "#595AB7",
    "yellowgreen": "#738A05"
}

// Generate test data， 每个节点随机连接到x个节点，x在0~Maxcon中随机。
let nodes = [];
let numNodes = 100;
let Maxcon = 10 //每个节点的
for (let x = 0; x < numNodes; x++) {
    let targetAry = [];
    let connections = (Math.round(Math.random() * Maxcon));
    for (let y = 0; y < connections; y++) {
        targetAry.push(Math.round(Math.random() * numNodes))
    }
    nodes.push({
        id: x,
        name: "Node_" + x,
        target: targetAry
    })
}

// Create the links array from the generated data
var links = [];
for (var i = 0; i < nodes.length; i++) {
    if (nodes[i].target !== undefined) {
        for (var j = 0; j < nodes[i].target.length; j++) {
            links.push({
                source: nodes[i],
                target: nodes[nodes[i].target[j]]
            })
        }
    }
}

try {
    request = new XMLHttpRequest();
}
catch (e) {
    try {
        request = new ActiveXObject();
    }
    catch (e) {
        alert("Your brower does not support XMLHTTP!");
    }
}

function geturl() {
    //http://39.104.200.8:18010/dump
    var url = "http://localhost:18010/dump";
    request.open("GET", url, false);
    request.send(null);
    var netdata1 = request.responseText
    var netdata = netdata1

    netdata = eval("(" + netdata + ")");
    console.log(netdata)
    dealdata(netdata)
}

// 保存node link 的name 和 id的对应关系
var map_node = new Map();
var map_link = new Map();

// 返回值 0表示节点不存在，非零表示已经存在
function pushNode(nodename) {
    if (map_node.has(nodename)) {

        return map_node.get(nodename)
    }
    map_node.set(nodename, nodeNum)
    nodeNum++
    return 0;
}

function pushLink(linkname) {
    if (map_link.has(linkname)) {
        // 如果link已存在返回零
        return 1
    } else {
        //如果未存在则设置其编号
        map_link.set(linkname, linkNum)
        linkNum++
    }
    return 0
}



// 每个记录是本结点加上本节点和别的节点的连接。优先填充本节点的信息，连接到的节点信息之后之后再补充。
function dealdata(netdata) {
    // 处理每个worker返回的结果。
    if (netdata[0].Reply === null) {
        // alert("dump 没有数据")
        console.log("dump没有数据")
        return
    }
    console.log("dump有数据。")
    for (var i = 0; i < netdata.length; i++) {
        var name_link = Object.keys(netdata[i].Reply.Links)
        console.log(name_link)
        let targetAry = [];
        var thisnodeid = nodeNum
        var thislink = netdata[i].Reply.Links
        var thisnodename = netdata[i].Reply.Node.ID
        var x = pushNode(thisnodename)

        if (x) {
            thisnodeid = x
        } else {
            mynode.push({
                id: thisnodeid,
                name: thisnodename,
            })
        }

        for (var j = 0; j < name_link.length; j++) {
            x = pushNode(name_link[j])
            if (x === 0) {
                mynode.push({
                    id: nodeNum - 1,
                    name: name_link[j],
                    target: []
                })
                targetAry.push(nodeNum)
            } else {
                targetAry.push(x)
            }
            x = pushLink(thisnodename + name_link[j])
            // 如果是新link，则压入mylink中
            if (x === 0) {
                mylink.push({
                    source: mynode[thisnodeid],
                    target: mynode[nodeNum - 1],
                    RateIn: thislink[name_link[j]].RateIn,
                    RateOut: thislink[name_link[j]].RateOut,
                    TotalIn: thislink[name_link[j]].TotalIn,
                    TotalOut: thislink[name_link[j]].TotalOut
                })
            }

        }
        mynode[thisnodeid].target = targetAry;
    }
    // 定义自己的node link
    console.log(mynode)
    console.log(mylink)
}

geturl()
// console.log(nodes)
// console.log(links)


nodes = mynode
links = mylink
numNodes = nodeNum

// Create SVG
var fdGraph = d3.select('#graphic svg')
// .attr('width', 0.8 * width)
// .attr('height', 0.8 * height)

// Create the force layout to calculate and animate node spacing
var forceLayout = d3.layout.force()
    .nodes(nodes)
    .links([])
    .gravity(gravity)
    .charge(charge)
    .size([0.6 * width, 0.8 * height])

// Create the SVG lines for the links
var link = fdGraph
    .selectAll('g').data(links).enter()
    .append('g')


link.append('line')
    .attr('x1', function (d) { return d.source.x })
    .attr('y1', function (d) { return d.source.y })
    .attr('x2', function (d) { return d.target.x })
    .attr('y2', function (d) { return d.target.y })
    .attr('stroke', palette.gray)
    .attr('stroke-width', 1)
    .attr('class', function (d, i) {
        // Add classes to lines to identify their from's and to's
        var theClass = 'from_' + d.source.id + ' line ';
        if (d.target !== undefined) {
            theClass += 'to_' + d.target.id
        }
        // line_1 line to_2 (ps:这里表示分别属于三个类，line1、line、to_2)
        return theClass
    })
    .on('mouseout', function (d) {
        d3.select(this).selectAll('text')
            .attr('font-size', '12')
            .attr('font-weight', 'normal')
        console.log('移动出line上')
    })
    .on('mouseover', function (d) {
        console.log(this)
        d3.select(this).selectAll('text')
            .attr('font-size', '16')
            .attr('font-weight', 'normal')
        console.log('移动到line')
    })


function nodeClick() {
    lock_click = 1;
    console.log("click node")
    // more Highlight the current node
    // debugger
    this.select('text')
        .attr('font-size', '18')
        .attr('font-weight', 'bold')
}

// Create the SVG groups for the nodes and their labels
var node = fdGraph
    .selectAll('circle').data(nodes).enter()
    .append('g')
    .attr('id', function (d) { return 'node_' + d.id })
    .attr('class', 'node')
    .on('mouseover', function (d) {
        // console.log(this)
        // When mousing over a node, make the label bigger and bold
        // and revert any previously enlarged text to normal
        if (lock_click === 1) return
        d3.selectAll('.node').selectAll('text')
            .attr('font-size', '12')
            .attr('font-weight', 'normal')

        // Highlight the current node
        d3.select(this).select('text')
            .attr('font-size', '16')
            .attr('font-weight', 'bold')

        // Hightlight the nodes that the current node connects to
        for (let i = 0; i < d.target.length; i++) {
            d3.select('#node_' + d.target[i]).select('text')
                .attr('font-size', '14')
                .attr('font-weight', 'bold')
        }

        // Reset and fade-out the unrelated links
        d3.selectAll('.line')
            .attr('stroke', palette.lightgray)
            .attr('stroke-width', 1)

        // link connect to this node is orange
        // link connect from this node to other is purple
        for (let x = 0; x < links.length; x++) {
            if (links[x].target !== undefined) {
                if (links[x].target.id === d.id) {
                    // Highlight the connections to this node
                    d3.selectAll('.to_' + d.id)
                        .attr('stroke', palette.orange)
                        .attr('stroke-width', 2)

                    // Highlight the nodes connected to this one
                    d3.select('#node_' + links[x].source.id).select('text')
                        .attr('font-size', '14')
                        .attr('font-weight', 'bold')
                }
            }
        }

        // Highlight the connections from this node
        d3.selectAll('.from_' + d.id)
            .attr('stroke', palette.purple)
            .attr('stroke-width', 3)

        // When mousing over a node, 
        // make it more repulsive so it stands out 让他排斥周围的节点而突出
        forceLayout.charge(function (d2, i) {
            if (d2 != d) {

                // Make the nodes connected to more repulsive
                for (let i = 0; i < d.target.length; i++) {
                    if (d2.id == d.target[i]) {
                        return charge * 8
                    }
                }

                // Make the nodes connected from more repulsive
                for (var x = 0; x < links.length; x++) {
                    if (links[x].source.id === d2.id) {
                        if (links[x].target !== undefined) {
                            if (links[x].target.id === d.id) {
                                return charge * 8
                            }
                        }
                    }
                }

                // Reset unrelated nodes
                return charge * 1;

            } else {
                // Make the selected node more repulsive
                return charge * 10;
            }
        });
        forceLayout.start();
    })
    .on("click", () => (nodeClick.call(this)))
    .call(forceLayout.drag)

// Create the SVG circles for the nodes
node.append('circle')
    .attr('cx', function (d) {
        return d.x
    })
    .attr('cy', function (d) {
        return d.y
    })
    .attr('r', circleWidth)
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

// Create the SVG text to label the nodes
node.append('text')
    .text(function (d) {
        return d.name.slice(0, 3) + "***" + d.name.slice(-3)
    })
    .attr('font-size', '12')

link
    .append('text')
    .attr('transform', function (d) {
        console.log("222")
        var a = d.source
        console.log(a, d)
        return `translate(` + 100 + ',' + 100 + ')'
    })
    // .text(function (d) {
    //     return "I am Text"
    //     //strconv.Itoa(TotalIn) + strconv.Itoa(TotalOut)
    // })
    // .attr('font-size', '12')

var nihao = true
var nihaonihao = 0
// Animate the layout every time tick
forceLayout.on('tick', function (e) {
    // Move the nodes base on charge and gravity
    node.attr('transform', function (d, i) {
        return 'translate(' + d.x + ', ' + d.y + ')'
    })

    // link.attr('transform', function (d, i) {
    //     return 'translate(' + d.x + ', ' + d.y + ')'
    // })

    // Adjust the lines to the new node positions
    link.selectAll('line')
        .attr('x1', function (d) {
            return d.source.x
        })
        .attr('y1', function (d) {
            return d.source.y
        })
        .attr('x2', function (d) {

            if (false) {
                console.log(d.target.x)
                nihaonihao++
                if (nihaonihao === 1)
                    nihao = false
            }

            if (d.target !== undefined) {
                return d.target.x
            } else {
                return d.source.x
            }
        })
        .attr('y2', function (d) {
            if (d.target !== undefined) {
                return d.target.y
            } else {
                return d.source.y
            }
        })
})

// Start the initial layout
forceLayout.start();


$(document).click(function (e) {
    let line_class = $('line');
    let node_class = $('circle');
    // console.log(e)
    // return
    if (!line_class.is(e.target) && line_class.has(e.target).length === 0) {
        if (!node_class.is(e.target) && node_class.has(e.target).length === 0) {
            d3.selectAll(".line")
                .attr('stroke', palette.lightgray)
                .attr('stroke-width', 1)
            d3.selectAll('.node').selectAll('text')
                .attr('font-size', '12')
                .attr('font-weight', 'normal')
            console.log("点击空白")
            lock_click = 0;

            //回复charge
            forceLayout.charge(function (d2, i) {
                return charge;
            });
            forceLayout.start();
        }
    }
});

d3.select(".searchBtn")
    .on("click", function (e) {
        let ee = document.getElementById("okk").value;
        console.log(ee)
        let el = document.querySelector("#node_" + ee + " circle")
        nodeClick.call(d3.select(el))
        // let x = document.getElementById("Node_1")// .dispatchEvent('mouseover');
        // console.log(x)
    })

d3.select('#graphic')
    .attr("float", 'left')

d3.select('#search')
    .attr("float", 'right')

function showHint(str) {
    if (str.length == 0) {
        return;
    }
    console.log(str, "触发了")
}