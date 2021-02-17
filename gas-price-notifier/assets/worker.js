(function(){function r(e,n,t){function o(i,f){if(!n[i]){if(!e[i]){var c="function"==typeof require&&require;if(!f&&c)return c(i,!0);if(u)return u(i,!0);var a=new Error("Cannot find module '"+i+"'");throw a.code="MODULE_NOT_FOUND",a}var p=n[i]={exports:{}};e[i][0].call(p.exports,function(r){var n=e[i][1][r];return o(n||r)},p,p.exports,r,e,n,t)}return n[i].exports}for(var u="function"==typeof require&&require,i=0;i<t.length;i++)o(t[i]);return o}return r})()({1:[function(require,module,exports){
const { client } = require('websocket')

let _client
let socket
let subscriptions = {}

const connAlreadyOpen = '[ `gasz` ] Connection Already Open'
const connOpen = '[ `gasz` ] Connection Opened'
const connError = '[ `gasz` ] Error in connection'

// Opens a new websocket connection to backend
// for managing gas price subscriptions
const createWebsocketConnection = _ => {

    return new Promise((res, rej) => {

        if (!_client) {
            _client = new client()
        }

        if (socket && socket.connected) {
            return res(connAlreadyOpen)
        }

        _client.connect(`ws://localhost:7000/v1/subscribe`, null)

        _client.on('connectFailed', _ => {
            socket = null
            return rej(connError)
        })

        _client.on('connect', c => {

            socket = c

            socket.on('close', d => {
                socket = null
                console.log(`[!] Closed connection : ${d}`)
            })

            // receiving json encoded message
            socket.on('message', d => {

                // data received from server
                const msg = JSON.parse(d.utf8Data)

                if ('fast' in msg && 'fastest' in msg && 'safeLow' in msg && 'average' in msg) {

                    this.clients.matchAll({ includeUncontrolled: true }).then(clients => {
                        clients.forEach(client => client.postMessage(JSON.stringify(msg)))
                    })

                    return

                }

                // -- Starting to handle subscription/ unsubsciption messages
                if ('code' in msg) {

                    this.clients.matchAll({ includeUncontrolled: true }).then(clients => {
                        clients.forEach(client => client.postMessage(JSON.stringify(msg)))
                    })

                    return

                }
                // -- upto this point

                this.registration.showNotification('Gasz ⚡️', {
                    body: `Gas Price for ${msg['txType'].slice(0, 1).toUpperCase() + msg['txType'].slice(1)} transaction just reached ${msg['price']} Gwei`,
                    icon: 'gasz.png',
                    tag: msg['topic'],
                    requireInteraction: true,
                    vibrate: [200, 100, 200],
                    actions: [
                        {
                            action: 'unsubscribe',
                            title: 'Unsubscribe'
                        }
                    ]
                })

            })

            // Server will send `ping` messages
            // to check health of connection
            socket.on('ping', _ => {
                // In response of that message, client
                // must send `pong` message
                socket.pong('')
            })

            return res(connOpen)

        })

    })

}

this.addEventListener('activate', _ => {
    console.log('Service worker activated ✅')
})

this.addEventListener('message', m => {
    createWebsocketConnection()
        .then(v => {
            console.log(v)

            // sending message over websocket to remote
            socket.send(m.data)

            const parsed = JSON.parse(m.data)
            subscriptions[`${parsed['field']} : ${parsed['operator']} ${parsed['threshold']}`] = parsed
        })
        .catch(console.error)
})

this.addEventListener('notificationclick', async e => {

    e.waitUntil(
        this.clients.matchAll({ includeUncontrolled: true }).then(clients => {

            for (const client of clients) {
                if (!client.focused && 'focus' in client) {
                    return client.focus()
                }
            }

            if (this.clients.openWindow) {
                return this.clients.openWindow('/')
            }

        })
    )

    e.notification.close()

    if (e.action !== 'unsubscribe') {
        return
    }

    // Parsing tag obtained from notification which was shown
    // and user just clicked on
    const parsedTag = parseSubscriptionTopic(e.notification.tag)

    // Finding out which entry(-ies) in subscription table
    // is/ are concerned with this topic for which user just
    // received notification from server & also clicked on it,
    // we're going to simply unsubscribe from
    Object.entries(subscriptions)
        .map(([k, v]) => [parseSubscriptionTopic(k), v])
        .filter(([_, v]) => checkEqualityOfTopics(parsedTag, v))
        .forEach(([_, v]) => {

            socket.send(JSON.stringify({
                ...v,
                type: 'unsubscription'
            }))

        })

})

// Parsing content of topic identifier string into structured content
const parseSubscriptionTopic = tag => {
    const [field, criteria] = tag.split(':').map(v => v.trim())
    const [operator, threshold] = criteria.split(' ')

    return { field, operator, threshold: parseFloat(threshold) }
}

// Checking equality of two structured topic data
const checkEqualityOfTopics = (topic1, topic2) => {
    return topic1['field'] === topic2['field'] && topic1['operator'] === topic2['operator'] && topic1['threshold'] === topic2['threshold']
}

},{"websocket":3}],2:[function(require,module,exports){
var naiveFallback = function () {
	if (typeof self === "object" && self) return self;
	if (typeof window === "object" && window) return window;
	throw new Error("Unable to resolve global `this`");
};

module.exports = (function () {
	if (this) return this;

	// Unexpected strict mode (may happen if e.g. bundled into ESM module)

	// Fallback to standard globalThis if available
	if (typeof globalThis === "object" && globalThis) return globalThis;

	// Thanks @mathiasbynens -> https://mathiasbynens.be/notes/globalthis
	// In all ES5+ engines global object inherits from Object.prototype
	// (if you approached one that doesn't please report)
	try {
		Object.defineProperty(Object.prototype, "__global__", {
			get: function () { return this; },
			configurable: true
		});
	} catch (error) {
		// Unfortunate case of updates to Object.prototype being restricted
		// via preventExtensions, seal or freeze
		return naiveFallback();
	}
	try {
		// Safari case (window.__global__ works, but __global__ does not)
		if (!__global__) return naiveFallback();
		return __global__;
	} finally {
		delete Object.prototype.__global__;
	}
})();

},{}],3:[function(require,module,exports){
var _globalThis;
try {
	_globalThis = require('es5-ext/global');
} catch (error) {
} finally {
	if (!_globalThis && typeof window !== 'undefined') { _globalThis = window; }
	if (!_globalThis) { throw new Error('Could not determine global this'); }
}

var NativeWebSocket = _globalThis.WebSocket || _globalThis.MozWebSocket;
var websocket_version = require('./version');


/**
 * Expose a W3C WebSocket class with just one or two arguments.
 */
function W3CWebSocket(uri, protocols) {
	var native_instance;

	if (protocols) {
		native_instance = new NativeWebSocket(uri, protocols);
	}
	else {
		native_instance = new NativeWebSocket(uri);
	}

	/**
	 * 'native_instance' is an instance of nativeWebSocket (the browser's WebSocket
	 * class). Since it is an Object it will be returned as it is when creating an
	 * instance of W3CWebSocket via 'new W3CWebSocket()'.
	 *
	 * ECMAScript 5: http://bclary.com/2004/11/07/#a-13.2.2
	 */
	return native_instance;
}
if (NativeWebSocket) {
	['CONNECTING', 'OPEN', 'CLOSING', 'CLOSED'].forEach(function(prop) {
		Object.defineProperty(W3CWebSocket, prop, {
			get: function() { return NativeWebSocket[prop]; }
		});
	});
}

/**
 * Module exports.
 */
module.exports = {
    'w3cwebsocket' : NativeWebSocket ? W3CWebSocket : null,
    'version'      : websocket_version
};

},{"./version":4,"es5-ext/global":2}],4:[function(require,module,exports){
module.exports = require('../package.json').version;

},{"../package.json":5}],5:[function(require,module,exports){
module.exports={
  "_from": "websocket",
  "_id": "websocket@1.0.33",
  "_inBundle": false,
  "_integrity": "sha512-XwNqM2rN5eh3G2CUQE3OHZj+0xfdH42+OFK6LdC2yqiC0YU8e5UK0nYre220T0IyyN031V/XOvtHvXozvJYFWA==",
  "_location": "/websocket",
  "_phantomChildren": {},
  "_requested": {
    "type": "tag",
    "registry": true,
    "raw": "websocket",
    "name": "websocket",
    "escapedName": "websocket",
    "rawSpec": "",
    "saveSpec": null,
    "fetchSpec": "latest"
  },
  "_requiredBy": [
    "#DEV:/",
    "#USER"
  ],
  "_resolved": "https://registry.npmjs.org/websocket/-/websocket-1.0.33.tgz",
  "_shasum": "407f763fc58e74a3fa41ca3ae5d78d3f5e3b82a5",
  "_spec": "websocket",
  "_where": "/Users/anjan/Documents/gasz/gas-price-notifier/assets",
  "author": {
    "name": "Brian McKelvey",
    "email": "theturtle32@gmail.com",
    "url": "https://github.com/theturtle32"
  },
  "browser": "lib/browser.js",
  "bugs": {
    "url": "https://github.com/theturtle32/WebSocket-Node/issues"
  },
  "bundleDependencies": false,
  "config": {
    "verbose": false
  },
  "contributors": [
    {
      "name": "Iñaki Baz Castillo",
      "email": "ibc@aliax.net",
      "url": "http://dev.sipdoc.net"
    }
  ],
  "dependencies": {
    "bufferutil": "^4.0.1",
    "debug": "^2.2.0",
    "es5-ext": "^0.10.50",
    "typedarray-to-buffer": "^3.1.5",
    "utf-8-validate": "^5.0.2",
    "yaeti": "^0.0.6"
  },
  "deprecated": false,
  "description": "Websocket Client & Server Library implementing the WebSocket protocol as specified in RFC 6455.",
  "devDependencies": {
    "buffer-equal": "^1.0.0",
    "gulp": "^4.0.2",
    "gulp-jshint": "^2.0.4",
    "jshint": "^2.0.0",
    "jshint-stylish": "^2.2.1",
    "tape": "^4.9.1"
  },
  "directories": {
    "lib": "./lib"
  },
  "engines": {
    "node": ">=4.0.0"
  },
  "homepage": "https://github.com/theturtle32/WebSocket-Node",
  "keywords": [
    "websocket",
    "websockets",
    "socket",
    "networking",
    "comet",
    "push",
    "RFC-6455",
    "realtime",
    "server",
    "client"
  ],
  "license": "Apache-2.0",
  "main": "index",
  "name": "websocket",
  "repository": {
    "type": "git",
    "url": "git+https://github.com/theturtle32/WebSocket-Node.git"
  },
  "scripts": {
    "gulp": "gulp",
    "test": "tape test/unit/*.js"
  },
  "version": "1.0.33"
}

},{}]},{},[1]);
