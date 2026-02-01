import React, { useState } from 'react'
import Phone from './PhoneEmail/Phone'
import Email from './PhoneEmail/Email'

const phoneEmailMap = {
  phone: Phone,
  email: Email,
}

const StepPhoneEmail = ({ onNext }) => {
  const [type, setType] = useState('phone');
  const Component = phoneEmailMap[type];

  return (
    <>

      <div className='flex justify-center items-center'>
        <div>
          <div className='flex gap-2 mb-4 items-center justify-end'>
            <button className={`w-12 h-12 flex items-center justify-center cursor-pointer rounded-lg ${type === 'phone' ? 'bg-[#0077ff]' : 'bg-[#262626]'}`} onClick={() => setType('phone')}>
              <img src="/images/phone-white.png" alt="phone" />
            </button>
            <button className={`w-12 h-12 flex items-center justify-center cursor-pointer rounded-lg ${type === 'email' ? 'bg-[#0077ff]' : 'bg-[#262626]'}`} onClick={() => setType('email')}>
              <img src="/images/mail-white.png" alt="email" />
            </button>
          </div>

          <Component onNext={onNext} />
        </div>
      </div>
    </>
  )
}

export default StepPhoneEmail