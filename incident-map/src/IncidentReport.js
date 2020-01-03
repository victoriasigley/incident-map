import React, { Component } from 'react';
import './App.css';
import 'leaflet/dist/leaflet.css';
import Dropzone from 'react-dropzone';
import { Map, Marker, Popup, TileLayer } from 'react-leaflet';
import axios from 'axios';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
} from '@material-ui/core';

class IncidentMap extends Component {
  constructor(props) {
    super(props);

    this.state = {
      incident: {
        address: {},
        apparatus: [],
        description: {},
        fire_department: {},
      },
      uploadDialogOpen: false,
      uploadFile: null,
    };
    this.onDrop = this.onDrop.bind(this);
    this.onUpload = this.onUpload.bind(this);
  }

  componentDidMount() {
    this.cancelToken = axios.CancelToken.source();

    this.loadIncident();
  }

  async loadIncident() {
    const url = `http://${window.location.hostname}:8080/incident`;
    try {
      // const response = await axios.get('/incident', { cancelToken: this.cancelToken.token });
      const response = await axios.get(url, { cancelToken: this.cancelToken.token });
      this.setState({ incident: response.data });
    } catch (error) {
      console.error(error);
    }
  }

  async onUpload() {
    const formData = new FormData();
    formData.append('file', this.state.uploadFile);
    const url = `http://${window.location.hostname}:8080/upload`;
    try {
      // const response = await axios.post('/upload', formData, { cancelToken: this.cancelToken.token });
      const response = await axios.post(url, formData, { cancelToken: this.cancelToken.token });
      console.log('onUpload');
      console.log(response.data);
      this.setState({ incident: response.data });
    } catch (error) {
      console.error(error);
    }
    this.setState({ uploadDialogOpen: false, uploadFile: null });
  }

  onDrop(files) {
    this.setState({ uploadFile: files[0] });
  }

  render() {
    const incident = { ...this.state.incident };

    const position = [incident.address ? (incident.address.latitude || 39) : 39, incident.address ? (incident.address.longitude || -98) : -98];

    const mapMarkers = incident.apparatus ? incident.apparatus.map(a =>
      (Object.keys(a.unit_status)).map(u =>
      <Marker
        position={[a.unit_status[u].latitude, a.unit_status[u].longitude]}
      >
        <Popup>
          {a.unit_id}<br/>
          {u}<br/>
          {a.unit_status[u].timestamp}
        </Popup>
      </Marker>
      )
    ) : null;
    console.log('render');
    console.log(position);
    console.log(mapMarkers);

    return (
    <div className="App">
      <header className="App-header">
        <p>{incident.fire_department ? incident.fire_department.name : ''}</p>
        <Button
          style={{ float: 'right', margin: '10px' }}
          variant='contained'
          onClick={() => this.setState({ uploadDialogOpen: true })}
        >
          Upload Incident File
        </Button>
      </header>
      <p>{incident.description ? incident.description.comments : ''}</p>
      <div style={{ height: '500px', width: '100%' }}>
        <Map center={position} zoom={13}>
          <TileLayer
             attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
             url='https://{s}.tile.osm.org/{z}/{x}/{y}.png'
           />
        {mapMarkers}
        </Map>
      </div>
      <Dialog open={this.state.uploadDialogOpen}>
        <DialogTitle>Upload Incident File</DialogTitle>
        <DialogContent>
          {this.state.uploadFile ?
            this.state.uploadFile.name
          :
          <Dropzone onDrop={this.onDrop} multiple={false}>
            {({ getRootProps, getInputProps }) => (
              <section>
                <div {...getRootProps()} style={{ height: '200px' }}>
                  <input {...getInputProps()} />
                  <p>Click here to browse or drag and drop a file into this space</p>
                </div>
              </section>
            )
            }
          </Dropzone>
          }
        </DialogContent>
        <DialogActions>
          <Button
            variant='contained'
            onClick={this.onUpload}
            disabled={!this.state.uploadFile}
          >
            Upload
          </Button>
          <Button
            variant='contained'
            onClick={() => this.setState({ uploadDialogOpen: false, uploadFile: null })}
          >
            Cancel
          </Button>
        </DialogActions>
     </Dialog>
    </div>
    );
  }
}

export default IncidentMap;
