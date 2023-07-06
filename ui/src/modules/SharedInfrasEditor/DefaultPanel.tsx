import React, { useCallback, useEffect, useState } from 'react'
import './index.css'
import { Alert, Button, Card, Form, ListGroup, ListGroupItem, Nav } from 'react-bootstrap'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

const getClassNameByExecution = (execution: any) => {
  if (execution?.status?.status === "SUCCESS") {
    return 'shared-infra-diagram__default-panel__execution--success'
  }

  if (execution?.status?.status === "RUNNING") {
    return 'shared-infra-diagram__default-panel__execution--running'
  }

  return 'shared-infra-diagram__default-panel__execution'
}

const DefaultPanel = ({ sharedInfra, onSave, goToView }: any) => {
  const [name, setName] = useState(sharedInfra?.name || '')
  const [description, setDescription] = useState(sharedInfra?.description || '')
  const [providerConfigRef, setProviderConfigRef] = useState<any>(sharedInfra?.providerConfigRef || '')

  const [providersConfig, setProvidersConfig] = useState<any>([])

  const getProvidersConfigs = useCallback(async () => {
    const sharedInfraRes = await fetch(`http://localhost:8080/providers-configs`)
    const sharedInfra = await sharedInfraRes.json()

    setProvidersConfig(sharedInfra.items)
  }, [])

  useEffect(() => {
    getProvidersConfigs()
  }, [])

  useEffect(() => {
    setName(sharedInfra?.name)
    setDescription(sharedInfra?.description)
    setProviderConfigRef(sharedInfra?.providerConfigRef)
  }, [sharedInfra])
  
  const handleCreate = () => {
    console.log(providerConfigRef)
    onSave({
      name,
      namespace: "default",
      description,
      providerConfigRef,
    })
  } 

  return (
    <div className='shared-infra-diagram__default-panel'>
      <div>
        {sharedInfra && <FontAwesomeIcon onClick={goToView} className='mb-2' style={{cursor: 'pointer'}} icon="arrow-left" />}
        <Card.Title>{sharedInfra?.name}</Card.Title>
        <Form.Group className='mb-3'>
          <Form.Label>Name</Form.Label>
          <Form.Control type="text" placeholder='Type the of yout shared infra' value={name} onChange={e => setName(e.target.value)} />
        </Form.Group>
        <Form.Group className="mb-3" controlId="exampleForm.ControlTextarea1" >
          <Form.Label>Description</Form.Label>
          <Form.Control as="textarea" rows={3} value={description} onChange={e => setDescription(e.target.value) } />
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Providers configs</Form.Label>
          <Form.Select value={providerConfigRef} onChange={e => setProviderConfigRef(e.target.value)}>
            <option value="" disabled>Select a provider config</option>
            {providersConfig.map((i: any) => (
              <option value={i}>{i.name}</option>
            ))}
          </Form.Select>
        </Form.Group>
        <div className="d-grid gap-2">
          <Button onClick={handleCreate}>Save shared infra</Button>
        </div>
      </div>
    </div>
  )
  
}

export default DefaultPanel