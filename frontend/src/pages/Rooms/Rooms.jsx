import React, { useEffect } from 'react'
import RoomCard from '../../components/RoomCard/RoomCard'
import AddRoomModel from '../../components/AddRoomModel/AddRoomModel'
import { useState } from 'react'
import { getAllRooms } from '../../http/index.jsx'

const Rooms = () => {

  const [showModel, setShowModel] = useState(false)
  const [rooms, setRooms] = useState([])

  useEffect(() => {
    const fetchRooms = async () => {
      const { data } = await getAllRooms();
      setRooms(data)
    }
    fetchRooms();
  }, [])

  function openModel() {
    setShowModel(true)
  }


  return (
    <>
      <div className='flex justify-between mx-auto max-w-6xl py-8'>
        <div className='flex items-center gap-8  relative pb-3 after:content-[""] after:absolute after:w-1/4 after:h-[4px] after:bg-blue-600 after:bottom-0'>
          <span className='font-bold text-xl'>All voices rooms</span>
          <div className='flex items-center gap-2 bg-[#262626] rounded-full px-4 py-2'>
            <img src="/images/search-icon.png" alt="search" />
            <input type="text" placeholder='Search' className='border-none outline-none bg-transparent text-white' />
          </div>
        </div>
        <div className=''>
          <button onClick={openModel} className='flex items-center bg-[#20bd5f] py-[5px] px-[20px] rounded-full gap-1 cursor-pointer hover:bg-[#0f6632] transition-colors'>
            <img src="/images/add-room-icon.png" alt="add" />
            <span className='font-bold text-white'>Create room</span>
          </button>
        </div>
      </div>

      <div className='mx-auto max-w-6xl grid grid-cols-4 gap-4'>
        {
          rooms.map((room) => (
            <RoomCard key={room.id} room={room} />
          ))
        }
      </div>
      {showModel && <AddRoomModel onClose={() => setShowModel(false)} />}
    </>
  )
}

export default Rooms