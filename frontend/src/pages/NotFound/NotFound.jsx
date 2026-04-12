import React from 'react'
import { Link } from 'react-router-dom'

const NotFound = () => {
    return (
        <div className='flex flex-col items-center justify-center min-h-[70vh] text-center'>
            <h1 className='text-6xl font-bold text-white mb-4'>404</h1>
            <h2 className='text-2xl font-semibold text-white mb-2'>Page Not Found</h2>
            <p className='text-gray-400 mb-6'>The page you are looking for does not exist.</p>
            <Link to="/"
                className="px-6 py-3 bg-[#0077ff] hover:bg-[#005bb5] transition-color rounded-full font-bold text-white"
            >
                Go to Home
            </Link>
        </div>
    )
}

export default NotFound