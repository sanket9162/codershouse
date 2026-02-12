import React, { useState } from 'react'
import Button from '../../../components/Button/Button'
import Card from '../../../components/Card/Card'
import TextInput from '../../../components/Textinput/TextInput'
import { sendOTP } from '../../../http/index'
import { useDispatch } from 'react-redux'
import { setOtp } from '../../../store/authSlice'

const Phone = ({ onNext }) => {
  const [phone, setPhone] = useState('')
  const dispatch = useDispatch()

  async function sumbit() {
    try {
      const res = await sendOTP(phone)
      dispatch(setOtp({ phone, hash: res.data.hash, expiresAt: res.data.expiresAt }))
      onNext()
    } catch (error) {
      console.log(error)
    }
  }

  return (
    <Card title="Enter Phone Number" icon="phone">
      <TextInput value={phone} onChange={(e) => setPhone(e.target.value)} />
      <div className='flex flex-col items-center justify-center'>
        <div className='mt-2'>
          <Button text="Next" onClick={sumbit} />
        </div>
        <p className='w-4/5 text-sm mt-4 text-gray-400'>By entering your number, you&apos;re agreeing to our Terms of Service and Privacy Policy</p>
      </div>
    </Card>
  )
}

export default Phone