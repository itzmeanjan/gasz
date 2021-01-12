let socket

this.addEventListener('install', e => {
    console.log('Install : ', e)
})

this.addEventListener('activate', e => {
    console.log('Activate : ', e)

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
