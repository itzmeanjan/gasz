const socket = new WebSocket(`${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/v1/subscribe`)

socket.onopen = _ => {
    console.log('[ `gasz` ] Connection Opened')
}

socket.onclose = _ => {
    console.log('[ `gasz` ] Connection Closed')
}

socket.onerror = _ => {
    console.log('[ `gasz` ] Error in connection')
    socket.close()
}

this.addEventListener('connect', e => {
    
    const port = e.ports[0]

    console.log(`New connection`)
    
    port.addEventListener('message', m => {

        console.log(`Received : ${m.data}`)

        port.postMessage(socket.readyState)
    })

    port.start()

})
