package searchutil

import "fmt"

func (s *Search) searchDefault(text string, i int) {
	if s.match(text) {
		s.output(text, fmt.Sprint(i))
	}
}

func (s *Search) searchDefaultInvert(text string, i int) {
	if !s.match(text) {
		s.output(text, fmt.Sprint(i))
	}
}

func (s *Search) searchInFileWithContext(text string, strNumber int) {
	if s.match(text) {
		if s.isPreCtx {
			for i := len(s.preCtxTextBuf) - 1; i >= 0; i-- {
				if s.preCtxTextBuf[i] != "" {
					s.output(s.preCtxTextBuf[i], s.preCtxStrNumBuf[i])
				}
			}
		}
		s.isPreCtx = false

		s.output(text, fmt.Sprint(strNumber))

		s.afterCtxCount = s.afterContext
	} else {
		s.isPreCtx = true
		for i := len(s.preCtxTextBuf) - 1; i > 0; i-- {
			s.preCtxTextBuf[i] = s.preCtxTextBuf[i-1]
			s.preCtxStrNumBuf[i] = s.preCtxStrNumBuf[i-1]
		}
		if len(s.preCtxTextBuf) > 0 {
			s.preCtxTextBuf[0] = text
			s.preCtxStrNumBuf[0] = fmt.Sprint(strNumber)
		}

		if s.afterCtxCount > 0 {
			s.output(text, fmt.Sprint(strNumber))
			s.afterCtxCount--
			s.isPreCtx = false
		}
	}
}

func (s *Search) searchInFileWithContextInvert(text string, strNumber int) {
	if !s.match(text) {
		if s.afterCtxCount <= 0 {
			if len(s.preCtxTextBuf) > 0 {
				lastIndex := len(s.preCtxTextBuf) - 1
				if s.preCtxTextBuf[lastIndex] != "" {
					s.output(s.preCtxTextBuf[lastIndex], s.preCtxStrNumBuf[lastIndex])
				}
				for i := len(s.preCtxTextBuf) - 1; i > 0; i-- {
					s.preCtxTextBuf[i] = s.preCtxTextBuf[i-1]
					s.preCtxStrNumBuf[i] = s.preCtxStrNumBuf[i-1]
				}
				s.preCtxTextBuf[0] = text
				s.preCtxStrNumBuf[0] = fmt.Sprint(strNumber)
			} else {
				s.output(text, fmt.Sprint(strNumber))
			}
		}
		s.afterCtxCount--
	} else {
		for i := len(s.preCtxTextBuf) - 1; i >= 0; i-- {
			s.preCtxTextBuf[i] = ""
			s.preCtxStrNumBuf[i] = ""
		}
		s.afterCtxCount = s.afterContext
	}
}
