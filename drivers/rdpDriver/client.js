/*
 * Copyright (c) 2015 Sylvain Peyrefitte
 *
 * This file is part of mstsc.js.
 *
 * mstsc.js is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */
//


(function() {

	TypeMessageRTM_index=0;
	
	TypeMessageRTM_connect=TypeMessageRTM_index++;
	TypeMessageRTM_bmp=    TypeMessageRTM_index++;
	TypeMessageRTM_ready=  TypeMessageRTM_index++;
	TypeMessageRTM_success=TypeMessageRTM_index++;
	TypeMessageRTM_close=  TypeMessageRTM_index++;
	TypeMessageRTM_error=  TypeMessageRTM_index++;
	TypeMessageRTM_mouse=  TypeMessageRTM_index++;
	TypeMessageRTM_scancode=  TypeMessageRTM_index++;
	TypeMessageRTM_wheel=     TypeMessageRTM_index++;
	
	/**
	 * Mouse button mapping
	 * @param button {integer} client button number
	 */
	function mouseButtonMap(button) {
		switch(button) {
		case 0:
			return 1;
		case 2:
			return 2;
		default:
			return 0;
		}
	};
	function parse(string, opts) {
  		if (opts === void 0) {
    		opts = {};
  		}
  		var encoding = {
    		chars: 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/',
    		bits: 6
  		};
  		// Build the character lookup table:
  		if (!encoding.codes) {
    		encoding.codes = {};

    		for (var i = 0; i < encoding.chars.length; ++i) {
      			encoding.codes[encoding.chars[i]] = i;
    		}
  		} // The string must have a whole number of bytes:


  		if (!opts.loose && string.length * encoding.bits & 7) {
    		//throw new SyntaxError('Invalid padding');
  		} // Count the padding bytes:


  		var end = string.length;

  		while (string[end - 1] === '=') {
    		--end; // If we get a whole number of bytes, there is too much padding:

    		if (!opts.loose && !((string.length - end) * encoding.bits & 7)) {
      			throw new SyntaxError('Invalid padding');
    		}
  		} // Allocate the output:


  		var out = new (opts.out || Uint8Array)(end * encoding.bits / 8 | 0); // Parse the data:
  		var bits = 0; // Number of bits currently in the buffer
  		var buffer = 0; // Bits waiting to be written out, MSB first
  		var written = 0; // Next byte to write

  		for (var _i = 0; _i < end; ++_i) {
    	// Read one character from the string:
    		var value = encoding.codes[string[_i]];

    		if (value === undefined) {
      			throw new SyntaxError('Invalid character ' + string[_i]);
    		} // Append the bits to the buffer:


    		buffer = buffer << encoding.bits | value;
   			bits += encoding.bits; // Write out some bits if the buffer has a byte's worth:

    		if (bits >= 8) {
      			bits -= 8;
      			out[written++] = 0xff & buffer >> bits;
    		}
  		} // Verify that we have received just enough bits:


  		if (bits >= encoding.bits || 0xff & buffer << 8 - bits) {
    		throw new SyntaxError('Unexpected end of data');
  		}

  		return out;
	};

	/**
	 * Mstsc client
	 * Input client connection (mouse and keyboard)
	 * bitmap processing
	 * @param canvas {canvas} rendering element
	 */
	function Client(canvas) {
		this.canvas = canvas;
		// create renderer
		this.render = new Mstsc.Canvas.create(this.canvas); 
		this.socket = null;
		this.activeSession = false;
		this.install();
	}
	
	Client.prototype = {
		install : function () {

			var self = this;
			// bind mouse move event
			this.canvas.addEventListener('mousemove', function (e) {
				if (!self.socket) return;
				console.log("mousemove");
				
				/*
				var offset = Mstsc.elementOffset(self.canvas);
				//self.socket.emit('mouse', e.clientX - offset.left, e.clientY - offset.top, 0, false);
*/
				var offset = Mstsc.elementOffset(self.canvas);
				self.socket.send(JSON.stringify({Type:TypeMessageRTM_mouse,Msg:{
					X:Math.floor((e.clientX - offset.left)),
					Y:Math.floor((e.clientY - offset.top)),
					Button:0,
					IsPressed:false
				}}));
				
				
				//e.preventDefault || !self.activeSession();
				//return false;
			},false);
			this.canvas.addEventListener('mousedown', function (e) {
				if (!self.socket) return;
				console.log("mousedown");
				
				/*
				var offset = Mstsc.elementOffset(self.canvas);
				self.socket.emit('mouse', e.clientX - offset.left, e.clientY - offset.top, mouseButtonMap(e.button), true);
				*/
				var offset = Mstsc.elementOffset(self.canvas);
				self.socket.send(JSON.stringify({Type:TypeMessageRTM_mouse,Msg:{
					X:Math.floor((e.clientX - offset.left)),
					Y:Math.floor((e.clientY - offset.top)),
					Button:mouseButtonMap(e.button),
					IsPressed:true
				}}));
				
				e.preventDefault();
				return false;
			});
			this.canvas.addEventListener('mouseup', function (e) {
				if (!self.socket || !self.activeSession) return;
				console.log("mouseup");
				
				/*
				var offset = Mstsc.elementOffset(self.canvas);
				self.socket.emit('mouse', e.clientX - offset.left, e.clientY - offset.top, mouseButtonMap(e.button), false);
				*/
				var offset = Mstsc.elementOffset(self.canvas);
				self.socket.send(JSON.stringify({Type:TypeMessageRTM_mouse,Msg:{
					X:Math.floor((e.clientX - offset.left)),
					Y:Math.floor((e.clientY - offset.top)),
					Button:mouseButtonMap(e.button),
					IsPressed:false
				}}));
				
				e.preventDefault();
				return false;
			});
			this.canvas.addEventListener('contextmenu', function (e) {
				if (!self.socket || !self.activeSession) return;
				console.log("contextmenu");
				/*
				var offset = Mstsc.elementOffset(self.canvas);
				self.socket.emit('mouse', e.clientX - offset.left, e.clientY - offset.top, mouseButtonMap(e.button), false);
				*/
				var offset = Mstsc.elementOffset(self.canvas);
				self.socket.send(JSON.stringify({Type:TypeMessageRTM_mouse,Msg:{
					X:Math.floor((e.clientX - offset.left)*1.02),
					Y:Math.floor((e.clientY - offset.top)*1.02),
					Button:mouseButtonMap(e.button),
					IsPressed:false
				}}));
				
				e.preventDefault();
				return false;
			});
			this.canvas.addEventListener('DOMMouseScroll', function (e) {
				if (!self.socket || !self.activeSession) return;
				console.log("DOMMouseScroll");

				/*
				var isHorizontal = false;
				var delta = e.detail;
				var step = Math.round(Math.abs(delta) * 15 / 8);
				
				var offset = Mstsc.elementOffset(self.canvas);
				self.socket.emit('wheel', e.clientX - offset.left, e.clientY - offset.top, step, delta > 0, isHorizontal);
				*/

				var isHorizontal = false;
				var delta = e.detail;
				var step = Math.round(Math.abs(delta) * 15 / 8);
				
				var offset = Mstsc.elementOffset(self.canvas);

				self.socket.send(JSON.stringify({Type:TypeMessageRTM_wheel,Msg:{
					X:Math.floor((e.clientX - offset.left)*1.03),
					Y:Math.floor((e.clientY - offset.top)*1.03),
					Step:step,
					IsNeg:delta>0,
					IsH:isHorizontal
				}}));	
				e.preventDefault();
				return false;
			});
			this.canvas.addEventListener('mousewheel', function (e) {
				if (!self.socket || !self.activeSession) return;
				console.log("mousewheel");
				/*
				var isHorizontal = Math.abs(e.deltaX) > Math.abs(e.deltaY);
				var delta = isHorizontal?e.deltaX:e.deltaY;
				var step = Math.round(Math.abs(delta) * 15 / 8);
				
				var offset = Mstsc.elementOffset(self.canvas);
				self.socket.emit('wheel', e.clientX - offset.left, e.clientY - offset.top, step, delta > 0, isHorizontal);
				*/

				var isHorizontal = Math.abs(e.deltaX) > Math.abs(e.deltaY);
				var delta = isHorizontal?e.deltaX:e.deltaY;
				var step = Math.round(Math.abs(delta) * 15 / 8);
				
				self.socket.send(JSON.stringify({Type:TypeMessageRTM_wheel,Msg:{
					X:Math.floor((e.clientX - offset.left)*1.03),
					Y:Math.floor((e.clientY - offset.top)*1.03),
					Step:step,
					IsNeg:delta>0,
					IsH:isHorizontal
				}}));
				e.preventDefault();
				return false;
			});
			
			// bind keyboard event
			window.addEventListener('keydown', function (e) {
				if (!self.socket || !self.activeSession) return;
				console.log("keydown");
				
				//self.socket.emit('scancode', Mstsc.scancode(e), true);
				self.socket.send(JSON.stringify({Type:TypeMessageRTM_scancode,Msg:{
					Button:Mstsc.scancode(e),
					IsPressed:true
				}}));
				e.preventDefault();
				return false;
			});
			window.addEventListener('keyup', function (e) {
				if (!self.socket || !self.activeSession) return;
				console.log("keyup");
				//self.socket.emit('scancode', Mstsc.scancode(e), false);
				
				self.socket.send(JSON.stringify({Type:TypeMessageRTM_scancode,Msg:{
					Button:Mstsc.scancode(e),
					IsPressed:false
				}}));
				
				e.preventDefault();
				return false;
			});
			
			return this;
		},
		/**
		 * connect
		 * @param ip {string} ip target for rdp
		 * @param domain {string} microsoft domain
		 * @param username {string} session username
		 * @param password {string} session password
		 * @param next {function} asynchrone end callback
		 */
		connect : function (ip, port, domain, username, password, next) {

			// start connection
			var self = this;
			var url = (window.location.protocol=="https"?"wss":"ws") + "://" + window.location.host+"/ws";
			console.log("connect to ",url);
			self.socket= new WebSocket(url);

			self.socket.onopen = function(event){
				console.log("connected webservice",event)
				self.socket.send(JSON.stringify({Type:TypeMessageRTM_connect,Msg:{
					Ip:ip,
					Port:port,
					Domain:domain,
					Username:username,
					Password:password,
					Width:$("#myCanvas").height(),
					Heigth:$("#myCanvas").width()
				}}));
			}
			
			self.socket.onmessage = function(msg){

				//console.log(msg.data)
				var dataProtocol=JSON.parse(msg.data)

				switch(dataProtocol.Type){
					default:
						console.log("Error Type data protocol:",dataProtocol.Type)
						break;
					case TypeMessageRTM_bmp:
						var arrayBMP=dataProtocol.Msg
						for (var i in arrayBMP)
						{ 
							var bitmap = arrayBMP[i];
							
							var raw = window.atob(arrayBMP[i].data);
							var rawLength = raw.length;
							var array = new Uint8Array(new ArrayBuffer(rawLength));

							for(i = 0; i < rawLength; i++) {
								array[i] = raw.charCodeAt(i);
							}
							bitmap["data"]=array;
							//console.log(bitmap);
							//console.log('[mstsc.js] bitmap update:' + bitmap.bitsPerPixel);
							self.render.update(bitmap);
						}
						break;
					case TypeMessageRTM_ready:
						console.log("RDP ready")
						self.activeSession = true;
						break;
				}
			}
			
			// 
			self.socket.onclose = function(){
			     console.log("Websocket: Closed")
			    setTimeout(function(){self.connect(ip,port,domain,username,password)}, 3000)
			}
			
			// Handle any errors that occur.
			self.socket.onerror = function(error) {
			    console.log('WebSocket Error: ' + error);
			};
			/*
			console.log(url, path)
			this.socket = io(url, { "path": path }).on('rdp-connect', function() {
				// this event can be occured twice (RDP protocol stack artefact)
				console.log('[mstsc.js] connected');
				self.activeSession = true;
			}).on('rdp-bitmap', function(bitmaps) {
				for (var i in bitmaps)
				{ 
					var bitmap = bitmaps[i];
					bitmap["data"] = parse(bitmaps[i].data)
					console.log(bitmap);
					console.log('[mstsc.js] bitmap update:' + bitmap.bitsPerPixel);
					self.render.update(bitmap);
				}
			}).on('rdp-close', function() {
				next(null);
				console.log('[mstsc.js] close');
				self.activeSession = false;
			}).on('rdp-error', function (err) {
				next(err);
				console.log(err)
				console.log('[mstsc.js] error : ' + err["code"] + '(' + err["message"] + ')');
				alert('error : ' + err["message"])
				self.activeSession = false;
			});
			console.log("password:");
			console.log(password);
			// emit infos event
			this.socket.emit('infos', {
				ip : ip, 
				port : port, 
				screen : { 
					width : this.canvas.width, 
					height : this.canvas.height 
				}, 
				domain : domain, 
				username : username, 
				password : password, 
				locale : Mstsc.locale()
			});
			*/
		}
	}
	
	Mstsc.client = {
		create : function (canvas) {
			return new Client(canvas);
		}
	}
})();
