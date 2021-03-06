package main

import "io"

/*
  UI behavior:
	o Clicking on a row whose log entry pertains to a packet:
		(i)   Highlights in red the cell that was originally clicked
		(ii)  Highlights in yellow all other rows pertaining to the same packet (same SeqNo)
		(iii) Highlights in orange all other rows pertaining to packets acknowledging this
		packet (AckNo is same as SeqNo of original packet)
		(iv)  Removes highlighting on rows that were previously highlighted using this
		procedure
	o Clicking the (left-most) cell of a row toggles a dark frame around the client-half of
	the row
	o Clicking the (right-most) cell of a row toggles a dark frame around the server-half of the row
	o Hovering over any row zooms on the row
 */

const (
	inspectorJavaScript =
		`
		jQuery(document).ready(function(){
			$('td[seqno].nonempty').click(onLeftClick);
			$('tr').mouseenter(hilightRow);
			$('tr').dblclick(toggleFoldRow);
			$('td.time').click(rotateMarkClientRow);
			$('td.time-abs').click(rotateMarkServerRow);
			$('td:has(div.tooltip)').mouseenter(showTooltip);
			$('td:has(div.tooltip)').mouseleave(hideTooltip);
			buildGraph();
		})
		function showTooltip() {
			var tt = $('div.tooltip', this);
			tt.css("display", "block");
		}
		function hideTooltip() {
			var tt = $('div.tooltip', this);
			tt.css("display", "none");
		}
		function onLeftClick(e) {
			var seqno = $(this).attr("seqno");
			if (_.isUndefined(seqno) || seqno == "") {
				return;
			}
			clearEmphasis();
			_.each($('[seqno='+seqno+'].nonempty'), function(t) { emphasize(t, "yellow-bkg") });
			_.each($('[ackno='+seqno+'].nonempty'), function(t) { emphasize(t, "orange-bkg") });
			emphasize($(this), "red-bkg");
		}
		function emphasize(t, bkg) {
			t = $(t);
			var saved_bkg = t.attr("emph");
			if (!_.isUndefined(saved_bkg)) {
				t.removeClass(saved_bkg);
			}
			t.addClass(bkg);
			t.attr("emph", bkg);
		}
		function clearEmphasis() {
			_.each($('[emph]'), function(t) {
				t = $(t);
				var saved_bkg = t.attr("emph");
				t.removeAttr("emph");
				t.removeClass(saved_bkg);
			});
		}
		function hilightRow() {
			_.each($('[hi]'), deHilightRow);
			$(this).attr("hi", 1);
			$('td', this).addClass("hi-bkg");
		}
		function toggleFoldRow() {
			$(this).addClass("folded");
		}
		function deHilightRow(t) {
			t = $(t);
			$('td', t).removeClass("hi-bkg");
			t.removeAttr("hi");
		}
		function rotateMarkClientRow() {
			var trow = $(this).parents()[0];
			_.each($('td.client, td.time', trow), _rotateMark);
		}
		function _rotateMark(t) {
			t = $(t);
			if (t.hasClass("mark-0")) {
				t.removeClass("mark-0");
				t.addClass("mark-1");
			} else if (t.hasClass("mark-1")) {
				t.removeClass("mark-1");
				t.addClass("mark-2");
			} else if (t.hasClass("mark-2")) {
				t.removeClass("mark-2");
			} else {
				t.addClass("mark-0");
			}
		}
		function rotateMarkServerRow() {
			var trow = $(this).parents()[0];
			_.each($('td.server, td.time-abs', trow), _rotateMark);
		}
		`
)

func printGraphJavaScript(w io.Writer, sweeper *SeriesSweeper) error {
	const one =
		`
		function buildGraph() {
			g = new Dygraph(
				document.getElementById("graph-box"),
		`
	const two =
		`
				, {
					connectSeparatedPoints: true,
					labelsDivWidth: 700,
					labelsDivStyles: { 'fontFamily': "Droid Sans Mono", 'fontWeight': "normal" },
					labelsSeparateLines: true,
					axes: {
						x: {
							axisLabelFormatter: function(x) {
								return x + 'ms';
							}
						}
					},
					drawPoints: true,
					colors: [ '#cc0000', '#00cc00', '#0000cc',' #00cccc', '#cc00cc', '#cccc00' ],
					labels: `
	const three =
		`
				}
			);
		}
		`
	w.Write([]byte(one))
	sweeper.EncodeData(w)
	w.Write([]byte(two))
	sweeper.EncodeHeader(w)
	_, err := w.Write([]byte(three))
	return err
}
