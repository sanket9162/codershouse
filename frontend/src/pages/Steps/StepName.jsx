import React, { useState } from 'react'
import Card from '../../components/Card/Card'
import TextInput from '../../components/TextInput/TextInput'
import Button from '../../components/Button/Button'
import { useDispatch, useSelector } from 'react-redux'
import { setName } from '../../store/activateSlice'

const StepName = ({ onNext }) => {
  const { name } = useSelector(state => state.activate)
  const dispatch = useDispatch()
  const [fullName, setFullName] = useState(name)
  function nextStep() {
    if (!fullName) return;
    onNext()

    dispatch(setName(fullName))

  }
  return (
    <>
      <div className='flex flex-col items-center justify-center'>

        <Card title="What's your full name?" icon="goggle-emoji">
          <TextInput value={fullName} onChange={(e) => setFullName(e.target.value)} />
          <div className='flex flex-col items-center justify-center'>
            <p className='w-4/5 text-sm mt-4 text-gray-400'>People use real names at codershouse !</p>

            <div className='mt-2'>
              <Button text="Next" onClick={nextStep} />
            </div>
          </div>
        </Card>

      </div>
    </>
  )
}

export default StepName