import React from 'react'

const RoomCard = ({ room }) => {
    return (
        <div className='bg-[#1d1d1d] p-[20px] rounded-xl cursor-pointer hover:bg-[#262626] transition-colors mb-5'>
            <h3 className=''>{room.topic}</h3>
            <div className='flex items-center gap-2 mt-4 justify-between'>
                <div className='flex my-auto'>
                    {room.speakers.map((speaker, index) => (
                        <img key={speaker.id} src={speaker.avatar} alt={speaker.name} width="40" height="40" className={`rounded-full object-cover border-2 bg-[#1d1d1d] border-blue-500 w-10 h-10 ${index !== 0 ? 'ml-[-15px] mt-[10px]' : ''}`} />
                    ))}
                </div>
                <div className=''>
                    {room.speakers.map((speaker) => (
                        <div key={speaker.id} className='flex items-center gap-2'>
                            <span>{speaker.name}</span>
                            <img src="/images/chat-bubble.png" alt="chat-bubble" />
                        </div>
                    ))}
                </div>
            </div>
            <div className='flex items-center gap-2 space-x-2 mt-4 place-self-end'>
                <span>{room.totalPeople}</span>
                <img src="/images/user-icon.png" alt="user" />
            </div>
        </div>
    )
}

export default RoomCard