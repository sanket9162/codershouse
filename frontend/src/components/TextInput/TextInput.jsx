import React from 'react'

const TextInput = (props) => {
    const { fullwidth, ...rest } = props;
    return (
        <div className={fullwidth === "true" ? "w-full" : ""}>
            <input type="text" {...rest} className={`mt-4 px-4 py-2 rounded-lg bg-[#262626] text-white border-none focus:outline-none focus:ring-2 focus:ring-[#0077ff] ${fullwidth ? 'w-full' : 'w-3/4'}`} />
        </div>
    )
}

export default TextInput