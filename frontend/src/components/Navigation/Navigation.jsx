import React from 'react'
import { Link } from 'react-router-dom'


const Navigation = () => {
  return (
    <nav className='mx-auto max-w-6xl py-8'>
        <Link to="/" className='flex items-center gap-2'>
            <img src="/images/logo.png" alt="logo" />
            <span className='font-bold text-xl'>Codershouse</span>
        </Link>
    </nav>
  )
}

export default Navigation