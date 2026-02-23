import React from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { logout } from '../../http';
import { useDispatch, useSelector } from 'react-redux';
import { logoutAuth } from '../../store/authSlice';



const Navigation = () => {

  const dispatch = useDispatch();
  const navigate = useNavigate()
  const { isAuth, user } = useSelector((state) => state.auth)

  async function logoutUser() {
    try {
      await logout()
    } catch (err) {
      console.error("logout api failed, but clear local state", err)
    } finally {
      dispatch(logoutAuth())
      navigate("/")
    }

  }

  return (
    <nav className='flex justify-between mx-auto max-w-6xl py-8'>
      <div>

        <Link to="/" className='flex items-center gap-2'>
          <img src="/images/logo.png" alt="logo" />
          <span className='font-bold text-xl'>Codershouse</span>
        </Link>
      </div>

      {isAuth && (
        <div className='flex items-center gap-4'>
          <span>{user?.name}</span>
          <Link to="/rooms">
            <img
              src={user?.avatar ? `${import.meta.env.VITE_API_URL}${user.avatar}` : '/images/monkey-avatar.png'}
              alt="avatar"
              width="40"
              height="40"
              className='rounded-full object-cover border-2 border-blue-500 w-10 h-10'
            />
          </Link>
          <button
            onClick={logoutUser}
            className='bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-full transition-colors'
          >Logout</button>
        </div>
      )}
    </nav>
  )
}

export default Navigation