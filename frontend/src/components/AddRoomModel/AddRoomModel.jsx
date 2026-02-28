import React from 'react'
import TextInput from '../TextInput/TextInput'

const AddRoomModel = ({ onClose }) => {
    return (
        <div className='fixed top-0 left-0 right-0 bottom-0 bg-black/60 flex items-center justify-center z-50'>
            <div className='bg-[#1d1d1d] p-8 rounded-2xl w-full max-w-md relative border border-[#262626] shadow-lg'>
                <button onClick={onClose} className='absolute top-4 right-4 text-gray-400 hover:text-white transition-colors'>
                    X
                </button>
                <div>
                    <h3 className='text-xl font-bold mb-4'>Enter the topic to be discussed</h3>
                    <TextInput placeholder='Topic' />
                </div>
                <div></div>
            </div>
        </div>
    )
}

export default AddRoomModel