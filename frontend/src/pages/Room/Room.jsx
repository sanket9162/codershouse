import React, { useState } from 'react'
import { useWebRTC } from '../../hooks/useWebRTC'
import { useParams } from 'react-router-dom'
import { useSelector } from 'react-redux'


const Room = () => {
    const { id: roomId } = useParams()
    const user = useSelector((state) => state.auth.user)
    const { clients, provideRef } = useWebRTC(roomId, user)
    return (
        <>
            <div>All Connected Clients</div>
            {
                clients.map((client) => (
                    <div key={client.id}>
                        <audio
                            ref={(instance) => provideRef(instance, client.id)}
                            controls autoPlay></audio>
                        <h4>{client.name}</h4>
                    </div>
                ))
            }
        </>
    )
}

export default Room