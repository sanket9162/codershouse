import React from 'react'

const Card = ({title, icon, children}) => {
  return (
     <div className="w-full max-w-md bg-[#1d1d1d] p-10 rounded-2xl text-center mt-4">
          <div className="flex items-center justify-center mb-4 gap-4">
              <img src={`/images/${icon}.png`} alt="logo" />
              <h1 className="text-2xl font-bold">{title}</h1>
          </div>
         {children}
          
      </div>
  )
}

export default Card