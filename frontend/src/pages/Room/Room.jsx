import React, { useEffect, useState } from 'react'
import { useWebRTC } from '../../hooks/useWebRTC'
import { useParams } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { useNavigate } from 'react-router-dom'
import { getRoomById } from '../../http'


const Room = () => {
    const { id: roomId } = useParams()
    const user = useSelector((state) => state.auth.user)
    const { clients, provideRef, handleMute } = useWebRTC(roomId, user)
    const navigate = useNavigate()

    const [room, setRoom] = useState(null)
    const [isMuted, setIsMuted] = useState(true)
    const handleManualLeave = () => {
        navigate('/rooms')
    }

    useEffect(() => {
        handleMute(isMuted, user.id);
    }, [isMuted])

    useEffect(() => {
        const fetchRoom = async () => {
            try {
                const response = await getRoomById(roomId)
                const room = response.data
                console.log(room)
                setRoom(room)
            } catch (error) {
                console.error('Error fetching room:', error)
            }
        }
        fetchRoom()
    }, [roomId])


    const handleMuteClick = (clientId) => {
        if (clientId !== user.id) return;
        setIsMuted(isMuted => !isMuted)
    }

    return (
        <>
            <div>
                <div className='flex justify-between mx-auto max-w-6xl py-8'>
                    <div className='flex items-center ml-16 relative pb-3 after:content-[""] after:absolute after:w-2/4 after:h-[4px] after:bg-blue-600 after:bottom-0'>
                        <button onClick={handleManualLeave} className='flex items-center gap-2 cursor-pointer bg-none'>

                            <img src="/images/arrow-left.png" alt="back" />
                            <span className='font-semibold'>All voice rooms</span>


                        </button>
                    </div>
                </div>

            </div >
            <div className='bg-[#1d1d1d] w-full min-h-[calc(100vh-220px)] rounded-t-4xl'>
                <div className='mx-16 pt-8 flex justify-between'>
                    <h2 className='font-semibold text-xl'>{room ? room.topic : 'Loading...'}</h2>
                    <div className='flex gap-3'>
                        <button className='flex items-center gap-2 bg-[#262626] rounded-full px-4 py-2 cursor-pointer hover:bg-[#333333] transition-colors'>
                            <img src="/images/palm.png" alt="palm.png" />
                        </button>
                        <button onClick={handleManualLeave} className='flex items-center gap-2 bg-[#262626] rounded-full px-4 py-2 cursor-pointer hover:bg-[#333333] transition-colors'>
                            <img src="/images/win.png" alt="win.png" />
                            <span>Leave quietly</span>
                        </button>
                    </div>
                </div>

                <div className='mx-16 mt-8 flex items-center flex-wrap gap-4'>
                    {
                        clients.map((client) => {
                            return (
                                <div key={client.id} className='flex items-center flex-col gap-4 '>

                                    <div className='relative border-4 border-[#5453e0] rounded-full'>
                                        <audio
                                            ref={(instance) => provideRef(instance, client.id)}
                                            autoPlay></audio>
                                        <img src={client.avatar} alt={client.name} className='rounded-full w-[75px] h-[75px]' />
                                        <button onClick={() => handleMuteClick(client.id)} className='bg-[#212121] rounded-full p-2 cursor-pointer absolute bottom-0 right-0 w-[30px] h-[30px] '>
                                            {
                                                client.muted ? (
                                                    <img src="/images/mic-mute.png" alt="mic-mute.png" />
                                                ) : (
                                                    <img src="/images/mic.png" alt="mic.png" />
                                                )
                                            }
                                        </button>
                                    </div>
                                    <h4 className='font-semibold text-sm'>{client.name}</h4>
                                </div>
                            )
                        })
                    }
                </div>

            </div>
        </>
    )
}

export default Room