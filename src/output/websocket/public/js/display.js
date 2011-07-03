var Display = function (canvasElement) {
  // Constants
  var ScreenWidth  = 256;
  var ScreenHeight = 192;

  var BytesPerLine      = ScreenWidth / 8; // =32
  var BytesPerLine_log2 = 5;               // =log2(BytesPerLine)

  var ScreenWidth_Attr      = ScreenWidth / 8;  // =32
  var ScreenWidth_Attr_log2 = 5;                // =log2(ScreenWidth_Attr)
  var ScreenHeight_Attr     = ScreenHeight / 8; // =24

  var ScreenBorderX = 32;
  var ScreenBorderY = 32;

  // Screen dimensions, including the border
  var TotalScreenWidth  = ScreenWidth + ScreenBorderX*2;
  var TotalScreenHeight = ScreenHeight + ScreenBorderY*2;

  var SCREEN_BASE_ADDR = 0x4000;
  var ATTR_BASE_ADDR   = 0x5800;

  var bitmap_unpack_table = new Array();

  var changedRegions = new Array();

  var palette = ["rgb(0, 0, 0)",
		 "rgb(0, 0, 192)",
		 "rgb(192, 0, 0)",
		 "rgb(192, 0, 192)",
		 "rgb(0, 192, 0)",
		 "rgb(0, 192, 192)",
		 "rgb(192, 192, 0)",
		 "rgb(192, 192, 192)",
		 "rgb(0, 0, 0)",
		 "rgb(0, 0, 255)",
		 "rgb(255, 0, 0)",
		 "rgb(255, 0, 255)",
		 "rgb(0, 255, 0)",
		 "rgb(0, 255, 255)",
		 "rgb(255,255,0)",
		 "rgb(255, 255, 255)"];

  var pixels = new Array();

  var context = canvasElement.getContext('2d');

  // Initialize bitmap unpack table
  for (var a = 0; a < (1 << 8); a++)
  {
    var bitmap_unpack_table_a = new Array();
    bitmap_unpack_table_a[0] = (a >> 7) & 1;
    bitmap_unpack_table_a[1] = (a >> 6) & 1;
    bitmap_unpack_table_a[2] = (a >> 5) & 1;
    bitmap_unpack_table_a[3] = (a >> 4) & 1;
    bitmap_unpack_table_a[4] = (a >> 3) & 1;
    bitmap_unpack_table_a[5] = (a >> 2) & 1;
    bitmap_unpack_table_a[6] = (a >> 1) & 1;
    bitmap_unpack_table_a[7] = (a >> 0) & 1;
    bitmap_unpack_table[a] = bitmap_unpack_table_a;
  }

  this.renderUnscaled = function (displayData) {
    var X0 = ScreenBorderX;
    var Y0 = ScreenBorderY;

    var screen_dirty = displayData.dirty;
    var screen_attr = displayData.attr;
    var screen_bitmap = displayData.bitmap;

    var attr_x, attr_y;

	for (attr_y = 0; attr_y < ScreenHeight_Attr; attr_y++)
	{
	  var dst_Y0 = Y0 + 8*attr_y;
	  var attr_wy = ScreenWidth_Attr * attr_y;

		for (attr_x = 0; attr_x < ScreenWidth_Attr; attr_x++)
	        {
			if (screen_dirty[attr_wy+attr_x])
		        {
			  var dst_X0 = X0 + 8*attr_x;
			  var y = 0;
			  var src_ofs = ((8 * attr_y) << BytesPerLine_log2) + attr_x;
			  var dst_ofs = TotalScreenWidth*(dst_Y0+y) + dst_X0;
				while (y < 8)
				{
				  // Paper is in the lower 4 bits, ink is in the higher 4 bits
				  var paperInk = screen_attr[src_ofs];
				  var paperInk_array = [paperInk & 0xf, (paperInk >> 4) & 0xf];
				  var value = screen_bitmap[src_ofs];
				  var unpacked_value = bitmap_unpack_table[value];

				  for (var x = 0; x < 8; x++)
				  {
				    var color = paperInk_array[unpacked_value[x]];
				    pixels[dst_ofs+x] = color;
				  }

				  y += 1;
				  src_ofs += BytesPerLine;
				  dst_ofs += TotalScreenWidth;
				}

				changedRegions.push({x:dst_X0, y:dst_Y0, w:8, h:8});
			}
		}
	}
  };

  this.render = function render(displayData)
  {
    changedRegions = [];
    this.renderUnscaled(displayData);
    var cx, cy;
    for (var i = 0; i < changedRegions.length; i++)
    {
      var r = changedRegions[i];
      var end_x = r.x + r.w;
      var end_y = r.y + r.h;
      for (var y = r.y; y < end_y; y++)
      {
	var wy = TotalScreenWidth * y;
	cy = y * 2;
	for (var x = r.x; x < end_x; x++)
	{
	  cx = x * 2;
	  context.fillStyle = palette[pixels[wy+x]];
	  context.fillRect(cx, cy, 2, 2);
	}
      }
    }
  };

};