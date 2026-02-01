import React from 'react'

const TextInput = (props) => {
    return (
        <div>
            <input type="text" {...props} className='w-3/4 mt-4 px-4 py-2 rounded-lg bg-[#262626] text-white border-none focus:outline-none focus:ring-2 focus:ring-[#0077ff]' />
        </div>
    )
}

export default TextInput