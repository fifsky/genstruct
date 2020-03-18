import React, {useState} from 'react';
import {Button, Card, Col, Input, Row, Select} from 'antd';
import {genApi} from "../service";
import {sync} from "../util";

const {TextArea} = Input;
const { Option } = Select;

export default () => {
  const [schema, setSchema] = useState('');
  const [struct, setStruct] = useState('');
  const [tags, setTags] = useState(["db", "json"]);

  const onChange = ({target: {value}}) => {
    setSchema(value);
  };

  const handleTags = (value) => {
    setTags(value);
  };

  const onConvert = () => {
    if (schema === "") {
      return
    }
    sync(async function () {
      const ret = await genApi({"table": schema, "tags": tags})
      setStruct(ret);
    })
  };

  return (
    <div style={{margin: '20px'}}>
      <Card
        title="MySQL Structure to Golang Struct"
        extra={
          <Button type="primary" onClick={onConvert}>
            Convert
          </Button>
        }
      >
        <Row style={{marginBottom:'24px'}}>
          <Col>
            <Select mode="tags" style={{width: '100%'}} defaultValue={tags} placeholder="Tags" onChange={handleTags}>
              {tags}
            </Select>
          </Col>
        </Row>
        <Row gutter={{xs: 8, sm: 16, md: 24, lg: 32}}>
          <Col span={12}>
            <TextArea
              onChange={onChange}
              placeholder="Input MySQL Create Table Syntax"
              autoSize={{minRows: 20, maxRows: 20}}
            />
          </Col>
          <Col span={12}>
            <TextArea
              value={struct}
              readOnly
              autoSize={{minRows: 20, maxRows: 20}}
            />
          </Col>
        </Row>
      </Card>
    </div>
  );
};
