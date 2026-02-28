import React, { useState } from 'react'
import TextInput from '../TextInput/TextInput'


const AddRoomModel = ({ onClose }) => {
    const [roomType, setRoomType] = useState('open')
    const [topic, setTopic] = useState('');

    return (
        <div className='fixed top-0 left-0 right-0 bottom-0 bg-black/60 flex items-center justify-center z-50'>
            <div className='bg-[#1d1d1d] p-8 rounded-2xl w-full max-w-md relative border border-[#262626] shadow-lg '>
                <button onClick={onClose} className='absolute top-4 right-4 text-gray-400'>
                    <img src="/images/close.png" alt="close" />
                </button>
                <div className='mb-4 border-b pb-4 border-[#262626] border-b-[2px] flex flex-col w-full'>
                    <h3 className='text-xl font-bold mb-4'>Enter the topic to be discussed</h3>
                    <TextInput fullwidth="true" placeholder='Topic' value={topic} onChange={(e) => setTopic(e.target.value)} />
                    <h2 className='text-xl font-bold mt-4 mb-2'>Room Type</h2>
                    <div className='flex justify-between gap-2 my-4'>
                        <div
                            onClick={() => setRoomType('open')}
                            className={`flex flex-col items-center gap-2 cursor-pointer transition-colors rounded-xl p-4 w-1/3 ${roomType === 'open' ? 'bg-[#262626]' : 'hover:bg-[#262626]'}`}
                        >
                            <img src="/images/globe.png" alt="public" />
                            <span>Public</span>
                        </div>
                        <div
                            onClick={() => setRoomType('social')}
                            className={`flex flex-col items-center gap-2 cursor-pointer transition-colors rounded-xl p-4 w-1/3 ${roomType === 'social' ? 'bg-[#262626]' : 'hover:bg-[#262626]'}`}
                        >
                            <img src="/images/social.png" alt="social" />
                            <span>Social</span>
                        </div>
                        <div
                            onClick={() => setRoomType('private')}
                            className={`flex flex-col items-center gap-2 cursor-pointer transition-colors rounded-xl p-4 w-1/3 ${roomType === 'private' ? 'bg-[#262626]' : 'hover:bg-[#262626]'}`}
                        >
                            <img src="/images/lock.png" alt="private" />
                            <span>Private</span>
                        </div>
                    </div>
                </div>
                <div className='flex items-center flex-col gap-4 mt-4'>
                    <h2 className='text-xl font-bold'>Start a room, open to everyone</h2>
                    <button className='flex items-center gap-2 cursor-pointer bg-[#20bd5f] hover:bg-[#0f6632] transition-colors rounded-full px-4 py-2'><img src="/images/celebration.png" alt="celebration" /> Let's Go</button>
                </div>
            </div>
        </div>
    )
}

export default AddRoomModel