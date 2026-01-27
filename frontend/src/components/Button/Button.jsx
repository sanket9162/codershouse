import React from 'react'

const Button = ({text, onClick}) => {
  return (
    <button onClick={onClick} className="bg-[#0077ff] flex items-center gap-2 px-4 py-2 rounded-full m-auto mt-4 font-medium cursor-pointer hover:bg-[#0077ff]/80 transition-all duration-200">
              <span>{text}</span>
              <img src="/images/arrow-forward.png" alt="arrow-forward" />
    </button>
  )
}

export default Button 