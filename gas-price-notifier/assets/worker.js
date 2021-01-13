let socket
let subscriptions = {}

// Opens a new websocket connection to backend
// for managing gas price subscriptions
const createWebsocketConnection = _ => {

    const connAlreadyOpen = '[ `gasz` ] Connection Already Open'
    const connOpen = '[ `gasz` ] Connection Opened'
    const connClosed = '[ `gasz` ] Connection Closed'
    const connError = '[ `gasz` ] Error in connection'

    return new Promise((res, rej) => {

        if (socket && socket.readyState === socket.OPEN) {
            return res(connAlreadyOpen)
        }

        socket = new WebSocket(`ws://localhost:7000/v1/subscribe`)

        // websocket connection is open now
        socket.onopen = _ => {
            return res(connOpen)
        }

        // connection with backend got closed
        socket.onclose = _ => {
            return rej(connClosed)
        }

        // due to some error encountered, closing connection with backend
        socket.onerror = _ => {
            socket.close()

            return rej(connError)
        }

        // Handling case when message being received from server
        socket.onmessage = e => {
            // data received from server
            const msg = JSON.parse(e.data)

            // -- Starting to handle subscription/ unsubsciption messages
            if ('code' in msg) {
                self.clients.matchAll({ includeUncontrolled: true }).then(clients => {
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
                vibrate: [200, 100, 200]
            })
        }

    })

}

this.addEventListener('activate', _ => {

    // Checking whether already connected via websocket or not
    //
    // if not, new connection will be created
    createWebsocketConnection()
        .then(console.log)
        .catch(console.error)

})

this.addEventListener('message', m => {
    createWebsocketConnection()
        .then(v => {
            console.log(v)
            socket.send(m.data)

            // Keeping track of which topic this client is subscribed to
            const parsed = JSON.parse(m.data)
            subscriptions[`${parsed['field']} : ${parsed['operator']} ${parsed['threshold']}`] = parsed
        })
        .catch(console.error)
})

this.addEventListener('notificationclick', e => {

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

    e.notification.close()

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
