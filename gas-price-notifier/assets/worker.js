let socket

this.addEventListener('activate', _ => {

    socket = new WebSocket(`ws://localhost:7000/v1/subscribe`)

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

})

this.addEventListener('message', m => {
    console.log(`Received from client : '${m.data}', ${typeof m.data}`)

    m.source.postMessage(m.data)
})
