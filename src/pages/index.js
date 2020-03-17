import React, { useState } from 'react';
import { Card, Button, Row, Col, Input } from 'antd';
import styles from './index.less';

const { TextArea } = Input;

export default () => {
  const [schema, setSchema] = useState('');
  const [struct, setStruct] = useState('');

  const onChange = ({ target: { value } }) => {
    setSchema(value);
  };

  const onConvert = () => {
    setStruct(schema);
  };

  return (
    <div style={{ margin: '20px' }}>
      <Card
        title="MySQL Structure to Golang Struct"
        extra={
          <Button type="primary" onClick={onConvert}>
            Convert
          </Button>
        }
      >
        <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }}>
          <Col span={12}>
            <TextArea
              onChange={onChange}
              placeholder="Input MySQL Create Table Syntax"
              autoSize={{ minRows: 20, maxRows: 20 }}
            />
          </Col>
          <Col span={12}>
            <TextArea
              value={struct}
              readOnly
              autoSize={{ minRows: 20, maxRows: 20 }}
            />
          </Col>
        </Row>
      </Card>
    </div>
  );
};
