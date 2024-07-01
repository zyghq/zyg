const baseUrl = "http://localhost:3000";
var config,
  isHidden = !0,
  pageWidth = window.innerWidth;

function getWidgetConfig() {
  // TODO: make an API call, get widget config
  // for now mock promise response
  // fetch widget config specifically for the wigetId.
  const config = {
    // allow_only_domains: false,
    domainsOnly: false,
    domains: null,
    bubblePosition: "right",
    headerColor: "#9370DB",
    profilePicture: null,
    // initial_message: "👋 Hey, ask me anything!",
    // user_message_color: "#FFFFFF",
    // user_message_bg_color: "#9370DB",
    // bot_message_color: "#131216",
    // bot_message_bg_color: "#EEEEF1",
    // bubble_color: "#EEEEF1",
    // bubble_position: "right",
    // profile_picture: null,
    // name: "Zyg Chat",
    // header_color: "#9370DB",
    // show_initial_message: true,
  };
  // attach local config if available
  // update this once we make an API call.
  const { zygConfig: localConfig } = window;
  if (localConfig) {
    Object.assign(config, localConfig);
  }
  return Promise.resolve(config);
}

function hideZW() {
  var t = document.getElementById("zyg-frame");
  (t.style.opacity = 0),
    (t.style.transform = "scale(0)"),
    (t.style.position = "fixed"),
    (document.getElementById("zyg-button").style.display = "block"),
    (isHidden = !0);
}

function showZW() {
  var t = document.getElementById("zyg-frame");
  ((t.style.opacity = 1),
  (t.style.transform = "scale(1)"),
  (t.style.position = "fixed"),
  pageWidth < 768) &&
    (document.getElementById("zyg-button").style.display = "none");
  isHidden = !1;
}

function onMessageHandler(evt) {
  console.log("origin:", evt.origin);
  console.log("source", evt.source);
  console.log("data", evt.data);
  if (evt.origin !== baseUrl) return;
  if (evt.data === "close") {
    hideZW();
  }
}

function handlePageWidthChange() {
  pageWidth = window.innerWidth;
  var t = document.getElementById("zyg-frame"),
    e =
      pageWidth > 768
        ? "width: 448px; height: 85vh; max-height: 820px;"
        : "width: 100%; height: 100%; max-height: 100%; min-height: 100%; left: 0px; right: 0px; bottom: 0px; top: 0px;",
    i =
      "right" === config.bubblePosition
        ? "right: 16px; left: unset; transform-origin: right bottom;"
        : "left: 16px; right: unset; transform-origin: left bottom;",
    o = isHidden
      ? "opacity: 0 !important; transform: scale(0) !important;"
      : "opacity: 1 !important; transform: scale(1) !important;";
  t.style.cssText =
    "box-shadow: rgba(150, 150, 150, 0.2) 0px 10px 30px 0px, rgba(150, 150, 150, 0.2) 0px 0px 0px 1px; overflow: hidden !important; border: none !important; display: block !important; z-index: 2147483645 !important; border-radius: 0.75rem; bottom: 96px; transition: scale 200ms ease-out 0ms, opacity 200ms ease-out 0ms; position: fixed !important;" +
    i +
    e +
    o;
}

function createZygWidget(config) {
  if (config.domainsOnly && config.domains) {
    const domains = config.domains;
    const d = window.location.hostname;
    if (!domains.includes(d)) {
      console.log("domain not allowed...");
      return;
    }
  }

  var authToken = config?.identity?.authToken;

  // create the iframe parent div container
  var frameContainer = document.createElement("div");
  frameContainer.setAttribute("id", "zyg-frame");
  // add styling
  var fcs =
      pageWidth > 768
        ? "width: 448px; height: 85vh; max-height: 820px"
        : "width: 100%; height: 100%; max-height: 100%; min-height: 100%; left: 0px; right: 0px; bottom: 0px; top: 0px;",
    bbp =
      config.bubblePosition && "right" === config.bubblePosition
        ? "right: 16px; left: unset; transform-origin: right bottom;"
        : "left: 16px; right: unset; transform-origin: left bottom;";
  frameContainer.style.cssText =
    "position: fixed !important; box-shadow: rgba(150, 150, 150, 0.2) 0px 10px 30px 0px, rgba(150, 150, 150, 0.2) 0px 0px 0px 1px; overflow: hidden !important; opacity: 0 !important; border: none !important; display: none !important; z-index: 2147483645 !important; border-radius: 0.75rem; bottom: 96px; transition: scale 200ms ease-out 0ms, opacity 200ms ease-out 0ms; transform: scale(0) !important;" +
    bbp +
    fcs;

  // create the iframe
  var iframe = document.createElement("iframe");
  iframe.setAttribute("id", "zyg-iframe"),
    iframe.setAttribute("title", "Zyg Widget"),
    iframe.setAttribute("src", baseUrl),
    iframe.setAttribute("frameborder", "0"),
    iframe.setAttribute("scrolling", "no"),
    iframe.setAttribute(
      "style",
      "border: 0px !important; width: 100% !important; height: 100% !important; display: block !important; opacity: 1 !important;"
    ),
    frameContainer.appendChild(iframe), // append the Iframe to the parent div container.
    document.body.appendChild(frameContainer);

  var popButton = document.createElement("div");
  popButton.setAttribute("id", "zyg-button");

  // add styling to the button
  var pbs = "background-color:" + config.headerColor + ";";
  (pbs += "position: fixed; bottom: 1rem;"),
    (pbs +=
      config.bubblePosition && "right" === config.bubblePosition
        ? "right: 16px; left: unset;"
        : "left: 16px; right: unset;"),
    (pbs +=
      "width: 50px; height: 50px; border-radius: 25px; box-shadow: rgba(0, 0, 0, 0.2) 0px 4px 8px 0px; cursor: pointer; z-index: 2147483645;"),
    (pbs +=
      "transition: transform 0.2s ease-in-out, opacity 0.2s ease-in-out; transform: scale(0); opacity: 0;"),
    (popButton.style.cssText = pbs);

  var buttonInnerSt =
    '<div style="display: flex; align-items: center; justify-content: center; width: 100%; height: 100%; z-index: 2147483646;">';
  config.profilePicture
    ? (buttonInnerSt +=
        '<img src="' +
        config.profilePicture +
        '" style="width: 100%; height: 100%; border-radius: 100px;" />')
    : (loadDotLottiePlayer(),
      (buttonInnerSt +=
        '<dotlottie-player src="https://lottie.host/f4cbf306-18d7-4e25-af7c-48ba3186b36a/hfU4BU0B4Q.json" background="transparent" speed="1" style="width: 80%; height: 80%;" loop autoplay></dotlottie-player>')),
    (buttonInnerSt += "</div>"),
    (popButton.innerHTML = buttonInnerSt),
    document.body.appendChild(popButton); // now that the button is ready, append it to the root body.

  setTimeout(function () {
    (popButton.style.opacity = 1),
      (popButton.style.transform = "scale(1)"),
      (frameContainer.style.display = "block");
  }, 1e3),
    popButton.addEventListener("click", function () {
      isHidden ? showZW() : hideZW();
    });

  // iframe
  const iface = {
    authenticate: () => {
      const message = {
        event: "authenticate",
        payload: {
          authToken,
        },
      };
      const messageStr = JSON.stringify(message);
      iframe.contentWindow.postMessage(messageStr, baseUrl);
    },
  };
  iframe.addEventListener("load", function () {
    window.addEventListener("zyg:ping", function () {
      const event = new CustomEvent("zyg:pong", {
        detail: iface,
      });
      window.dispatchEvent(event);
    });
  });
}

async function loadDotLottiePlayer() {
  try {
    // @sanchitrk: move this to self hosted cdn.
    await import(
      "https://unpkg.com/@dotlottie/player-component@2.3.0/dist/dotlottie-player.mjs"
    );
  } catch (t) {
    console.error("Failed to load the DotLottie Player module", t);
  }
}

function widgetPing() {
  return new Promise((resolve) => {
    let timeoutId;
    const widgetLoadedBeat = () => {
      const event = new Event("zyg:ping");
      timeoutId = window.setTimeout(widgetLoadedBeat, 1000);
      window.dispatchEvent(event);
    };

    window.addEventListener(
      "zyg:pong",
      (e) => {
        const api = e.detail;
        window.clearTimeout(timeoutId);
        resolve(api);
      },
      {
        once: true,
      }
    );
    widgetLoadedBeat();
  });
}

function init() {
  getWidgetConfig()
    .then((c) => {
      console.log("widget config:", c);
      createZygWidget((config = c)),
        window.addEventListener("message", onMessageHandler),
        window.addEventListener("resize", handlePageWidthChange),
        widgetPing()
          .then((api) => {
            api.authenticate();
          })
          .catch((err) => {
            console.error("error on widget ping-pong", err);
          });
    })
    .catch((err) => {
      console.error("error fetching widget config:", err);
    });
}
init();
