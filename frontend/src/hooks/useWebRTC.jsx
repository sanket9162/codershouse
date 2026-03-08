import { useCallback, useEffect, useReducer, useRef } from "react";
import { useStateWithCallback } from "./useStateWithCallback";
import { socketInit } from "../socket/index"
import { ACTIONS } from "../action";



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
                socket.current.emit(ACTIONS.JOIN, { roomId, peerId: user.id })
            })
        })
    }, [])

    useEffect(() => {

        const handleNewPeer = async ({ peerId, createOffer, user: peerUser }) => {
            if (peerId in connections.current) {
                return console.log("already connected", peerId, user.name)
            }

            connections.current[peerId] = new RTCPeerConnection({
                iceServers: [
                    {
                        urls: [
                            'stun:stun.l.google.com:19302',
                            'stun:stun1.l.google.com:19302',
                        ],
                    },
                ],
            })

            // handle new ice candidate
            connections.current[peerId].onicecandidate = (event) => {
                if (event.candidate) {
                    socket.current.emit(ACTIONS.RELAY_ICE, {
                        peerId,
                        iceCandidate: event.candidate,
                    })
                }
            }


            // handle new remote stream
            connections.current[peerId].ontrack = ({ streams: [remoteStream] }) => {
                addNewClient(peerUser, () => {
                    const remoteAudio = audioElements.current[peerId]
                    if (!remoteAudio) {
                        return
                    }
                    remoteAudio.srcObject = remoteStream
                })
            }

            localMediaStream.current.getTracks().forEach(track => {
                connections.current[peerId].addTrack(
                    track,
                    localMediaStream.current
                )
            });

            if (createOffer) {
                const offer = await connections.current[peerId].createOffer()
                await connections.current[peerId].setLocalDescription(offer)
                socket.current.emit(ACTIONS.RELAY_SDP, {
                    peerId,
                    sessionDescription: offer,
                })
            }

        }


        socket.current.on(ACTIONS.ADD_PEER, handleNewPeer)

        return () => {
            socket.current.off(ACTIONS.ADD_PEER)
        }

    }, []);

    return { clients, provideRef }
}