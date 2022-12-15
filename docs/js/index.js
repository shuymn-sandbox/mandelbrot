const go = new Go();

const LIMIT = 100;
const CANVAS_HEIGHT = 900;
const CANVAS_WIDTH = 900;
const OFFSET_WIDTH_DEFAULT = 1.5;
const OFFSET_HEIGHT_DEFAULT = 1.0;
const SCALE_DEFAULT = 1.0;

/** @type {(x: number, y: number) => [[string, boolean][]]} */
let calculate;

WebAssembly.instantiateStreaming(fetch("wasm/main.wasm"), go.importObject).then(
  (obj) => {
    const wasm = obj.instance;
    go.run(wasm);

    calculate = global.calculate;
  }
);

/**
 * @param {number} ms
 */
async function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function setup() {
  createCanvas(CANVAS_WIDTH, CANVAS_HEIGHT);
  background(255);
  fill(32);
  noStroke();
  noLoop();

  for (let count = 0; count < 10; count++) {
    if (calculate != null) {
      console.debug(`wait ${count} times.`);
      break;
    }
    sleep(10 + 2 ** count);
  }
}

function draw() {
  let scale = SCALE_DEFAULT;
  const paramScale = getParam("scale");
  if (paramScale !== null) {
    scale = parseFloat(paramScale);
  }

  let offsetWidth = OFFSET_WIDTH_DEFAULT;
  const paramOffsetWidth = getParam("offsetWidth");
  if (paramOffsetWidth !== null) {
    offsetWidth = parseFloat(paramOffsetWidth);
  }

  let offsetHeight = OFFSET_HEIGHT_DEFAULT;
  const paramOffsetHeight = getParam("offsetHeight");
  if (paramOffsetHeight !== null) {
    offsetHeight = parseFloat(paramOffsetHeight);
  }

  calculate(
    LIMIT,
    scale,
    CANVAS_HEIGHT,
    CANVAS_WIDTH,
    offsetWidth,
    offsetHeight
  ).forEach((rows, x) => {
    rows.forEach(([color, _], y) => {
      stroke(color);
      point(x, y);
    });
  });

  console.debug("finished.");
}

/**
 * @param {string} key
 */
function getParam(key) {
  const url = new URL(window.location.href);
  const params = new URLSearchParams(url.search);
  return params.get(key);
}
