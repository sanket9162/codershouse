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

        return () => {
            localMediaStream.current?.getTracks().forEach((track) => track.stop())
            socket.current.emit(ACTIONS.LEAVE, { roomId })
            if (socket.current && socket.current.connected) {
                socket.current.disconnect()
            }
        }
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
        } else {
            // If already present, DOM is mounted, so execute callback safely immediately
            if (cb && typeof cb === "function") cb()
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
                if (localAudio) {
                    localAudio.volume = 0
                    localAudio.srcObject = localMediaStream.current
                }

                // socket emit JOIN socket io
                socket.current.emit(ACTIONS.JOIN, { roomId, peerId: user.id, ...user })
            })
        })
    }, [])

    useEffect(() => {

        const handleNewPeer = async ({ peerId, createOffer, user: peerUser }) => {
            if (peerId in connections.current) {
                return console.log("already connected", peerId, user.name)
            }

            addNewClient(peerUser, () => { })

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
                    const remoteAudio = audioElements.current[peerUser.id]
                    if (remoteAudio) {
                        remoteAudio.srcObject = remoteStream
                    } else {
                        // Fallback mapping in case DOM React 18 batches ref attachments
                        let settled = false
                        let interval = setInterval(() => {
                            if (audioElements.current[peerUser.id]) {
                                audioElements.current[peerUser.id].srcObject = remoteStream
                                settled = true
                                clearInterval(interval)
                            }
                        }, 50)
                        setTimeout(() => { if (!settled) clearInterval(interval) }, 2000)
                    }
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

    useEffect(() => {
        const handleIceCandidate = async ({ peerId, iceCandidate }) => {
            const connection = connections.current[peerId]
            if (connection) {
                if (!connection.remoteDescription) {
                    connection.iceCandidatesQueue = connection.iceCandidatesQueue || []
                    connection.iceCandidatesQueue.push(iceCandidate)
                } else {
                    try {
                        // RTCPeerConnection API requires remote description to be set before adding ICE candidates
                        await connection.addIceCandidate(new RTCIceCandidate(iceCandidate))
                    } catch (e) {
                        console.error("Error adding ice candidate", e)
                    }
                }
            }
        }

        socket.current.on(ACTIONS.ICE_CANDIDATE, handleIceCandidate)

        return () => {
            socket.current.off(ACTIONS.ICE_CANDIDATE)
        }
    }, []);


    // Handle SDP
    useEffect(() => {
        const handleRelaySDP = async ({ peerId, sessionDescription: remoteOffer }) => {
            const connection = connections.current[peerId]
            if (!connection) {
                return
            }
            await connection.setRemoteDescription(new RTCSessionDescription(remoteOffer))

            // Process any queued ICE candidates now that remote description is set!
            if (connection.iceCandidatesQueue) {
                connection.iceCandidatesQueue.forEach(async (candidate) => {
                    try {
                        await connection.addIceCandidate(new RTCIceCandidate(candidate))
                    } catch (e) {
                        console.error("Error queueing ice", e)
                    }
                })
                connection.iceCandidatesQueue = []
            }

            if (remoteOffer.type === 'offer') {
                const answer = await connection.createAnswer()
                await connection.setLocalDescription(answer)
                socket.current.emit(ACTIONS.RELAY_SDP, {
                    peerId,
                    sessionDescription: answer,
                })
            }
        }

        socket.current.on(ACTIONS.SESSION_DESCRIPTION, handleRelaySDP)

        return () => {
            socket.current.off(ACTIONS.SESSION_DESCRIPTION)
        }
    }, [])

    // handle remove peer
    useEffect(() => {
        const handleRemovePeer = ({ peerId, userId }) => {
            if (connections.current[peerId]) {
                connections.current[peerId].close()
                delete connections.current[peerId]
            }
            delete audioElements.current[userId]
            setClients((existingClients) => existingClients.filter((client) => client.id !== userId))
        }

        socket.current.on(ACTIONS.REMOVE_PEER, handleRemovePeer)

        return () => {
            socket.current.off(ACTIONS.REMOVE_PEER)
        }
    }, [])

    return { clients, provideRef }
}