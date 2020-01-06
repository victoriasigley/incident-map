import React, { Component } from 'react';
import './App.css';
import Dropzone from 'react-dropzone';
import { Map, Marker, Popup, TileLayer } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import moment from 'moment';

import axios from 'axios';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
} from '@material-ui/core';

delete L.Icon.Default.prototype._getIconUrl;

L.Icon.Default.mergeOptions({
      iconRetinaUrl: require('leaflet/dist/images/marker-icon-2x.png'),
      iconUrl: require('leaflet/dist/images/marker-icon.png'),
      shadowUrl: require('leaflet/dist/images/marker-shadow.png')
});

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
    const address = incident.address || {};
    const description = incident.description || {};
    const apparatus = incident.apparatus || [];

    const position = [(address.latitude || 39), (address.longitude || -98)];

    const mapMarkers = apparatus ? apparatus.map(a => {
      const filteredKeys = Object.keys(a.unit_status).filter(x => a.unit_status[x].latitude || a.unit_status[x].longitude);
      const sortedKeys = filteredKeys.sort((x, y) => {
        const isBefore = moment(a.unit_status[x].timestamp).isBefore(a.unit_status[y].timestamp);
        return isBefore ? 1 : -1;
      });
      return (sortedKeys.map((u, index) => {
        const isMostRecent = index === 0;
        const hasMarker = a.unit_status[u].latitude || a.unit_status[u].longitude;
        return (
          hasMarker ?
          <Marker
            key={index}
            position={[a.unit_status[u].latitude, a.unit_status[u].longitude]}
            opacity={isMostRecent ? 1 : 0.3}
          >
            <Popup>
              <b>{a.unit_id}</b><br/>
              {u}<br/>
              {moment(a.unit_status[u].timestamp).format("MMM D YYYY, h:mm:ss a")}
          </Popup>
        </Marker>
        : null
      );
      }));
    }
    ) : null;

    return (
    <div className="App">
      <header className="App-header">
        {incident.fire_department ?
        <p>
          {incident.fire_department.name}<br/>
          Shift {incident.fire_department.shift}
        </p>
        : ''}
        <Button
          style={{ float: 'right', margin: '10px' }}
          variant='contained'
          onClick={() => this.setState({ uploadDialogOpen: true })}
        >
          Upload Incident File
        </Button>
      </header>
      <p>{description.comments}</p>
      <div style={{ height: '500px', width: '100%' }}>
        <Map center={position} zoom={13}>
          <TileLayer
             attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
             url='https://{s}.tile.osm.org/{z}/{x}/{y}.png'
           />
           {address.latitude && address.longitude &&
           <Marker position={position}>
             <Popup>
               <b><u>INCIDENT LOCATION</u></b><br/>
               <b>Address:</b><br/>
               {address.common_place_name}<br/>
               {address.address_line1}<br/>
               {address.city}, {address.state} {address.postal_code}
               <br/><br/>
               <b>Cross Street: </b>{address.cross_street1} AND {address.cross_street2}
             </Popup>
           </Marker>
           }
           {mapMarkers}
         </Map>
         {description.extended_data &&
         <p>
           <b>Dispatch Statistics:</b><br/>
           Dispatch Duration: {description.extended_data.dispatch_duration} seconds<br/>
           Event Duration: {description.extended_data.event_duration} seconds<br/>
           Response Time: {description.extended_data.response_time} seconds<br/>
         </p>
         }
         {description.first_unit_dispatched &&
           <p>
             <b>First unit dispatched:</b> {moment(description.first_unit_dispatched).format("MMM D YYYY, h:mm:ss a")}
          </p>
          }
         {description.first_unit_enroute &&
           <p>
             <b>First unit enroute:</b> {moment(description.first_unit_enroute).format("MMM D YYYY, h:mm:ss a")}
          </p>
          }
         {description.first_unit_arrived &&
           <p>
             <b>First unit arrived:</b> {moment(description.first_unit_arrived).format("MMM D YYYY, h:mm:ss a")}
          </p>
          }
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
