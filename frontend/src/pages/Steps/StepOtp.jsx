import React, { useState } from 'react'
import Card from '../../components/Card/Card'
import TextInput from '../../components/Textinput/TextInput'
import Button from '../../components/Button/Button'

const StepOtp = ({ onNext }) => {
  const [otp, setOtp] = useState('')
  return (
    <>
      <div className='flex flex-col items-center justify-center'>
        <Card title="Enter OTP" icon="lock-emoji">
          <TextInput value={otp} onChange={(e) => setOtp(e.target.value)} />
          <div className='flex flex-col items-center justify-center'>
            <div className='mt-2'>
              <Button text="Next" />
            </div>
            <p className='w-4/5 text-sm mt-4 text-gray-400'>By entering your number, you&apos;re agreeing to our Terms of Service and Privacy Policy</p>
          </div>
        </Card>
      </div>
    </>
  )
}

export default StepOtp