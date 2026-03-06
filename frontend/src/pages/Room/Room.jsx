import React, { useState } from 'react'
import { useWebRTC } from '../../hooks/useWebRTC'

const Room = () => {
    const { clients } = useWebRTC()
    return (
        <>
            <div>All Connected Clients</div>
            {
                clients.map((client) => (
                    <div key={client.id}>
                        <audio controls autoPlay></audio>
                        <h4>{client.name}</h4>
                    </div>
                ))
            }
        </>
    )
}

export default Room