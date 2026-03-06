import { useCallback, useEffect, useReducer, useRef } from "react";
import { useStateWithCallback } from "./useStateWithCallback";
import { socketInit } from "../socket/index"



export const useWebRTC = (roomId, user) => {
    const [clients, setClients] = useStateWithCallback([])
    const audioElements = useRef({})
    const connections = useRef({})
    const localMediaStream = useRef(null)
    const socket = useRef(null)

    useEffect(() => {
        socket.current = socketInit()

    }, [])

    const provideRef = (instance, userId) => {
        audioElements.current[userId] = instance
    }

    const addNewClient = useCallback((newClient, cb) => {
        const lookingFor = clients.find((client) => client.id === newClient.id)
        if (!lookingFor) {
            setClients((existingClients) => {
                const isAlreadyPresent = existingClients.find((client) => client.id === newClient.id)
                if (!isAlreadyPresent) {
                    return [...existingClients, newClient]
                }
                return existingClients
            }, cb)
        }
    }, [clients, setClients])

    useEffect(() => {
        const startMedia = async () => {
            try {
                const stream = await navigator.mediaDevices.getUserMedia({
                    audio: true,
                })
                localMediaStream.current = stream
            } catch (error) {
                console.error('Error accessing media devices.', error)
            }
        }
        startMedia().then(() => {
            addNewClient(user, () => {
                const localAudio = audioElements.current[user.id]
                if (!localAudio) {
                    return
                }
                localAudio.volume = 0
                localAudio.srcObject = localMediaStream.current

                // socket emit JOIN socket io
                socket.current.emit("join", {})
            })
        })
    }, [])

    return { clients, provideRef }
}