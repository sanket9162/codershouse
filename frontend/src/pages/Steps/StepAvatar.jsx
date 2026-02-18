import React, { useState } from 'react'
import Button from '../../components/Button/Button'
import Card from '../../components/Card/Card'
import { useSelector, useDispatch } from 'react-redux'
import { setAvatar } from '../../store/activateSlice'
import { activateUser } from '../../http'

const StepAvatar = ({ onNext }) => {
  const dispatch = useDispatch();
  const { name, avatar } = useSelector(state => state.activate)
  const [image, setImage] = useState('/images/monkey-avatar.png')
  function captureImage(e) {
    const file = e.target.files[0];
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onloadend = () => {
      setImage(reader.result);
      dispatch(setAvatar(reader.result));

    }
  }

  async function submit() {
    try {
      const response = await activateUser(name, avatar);
      console.log(response);
    } catch (err) {
      console.log(err);
    }
  }
  return (
    <>
      <div className='flex flex-col items-center justify-center'>

        <Card title={`Ok! ${name}`} icon="monkey-emoji">
          <p className='text-center text-gray-400 mb-4'>How's this photo ?</p>
          <div className='w-[110px] h-[110px] mx-auto overflow-hidden border-4 border-[#0077ff] rounded-full flex items-center justify-center'>
            <img src={image} alt="avatar" className='w-full h-full object-cover' />
          </div>

          <div>
            <input id="avatarInput" type="file" className='hidden' accept='image/*' onChange={captureImage} />
            <label htmlFor="avatarInput" className='text-[#0077ff] mt-2 inline-block cursor-pointer'>
              Choose a different photo
            </label>
          </div>

          <div className='mt-2'>
            <Button text="Next" onClick={submit} />
          </div>
        </Card>

      </div>
    </>
  )
}

export default StepAvatar