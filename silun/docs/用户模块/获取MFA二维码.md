# 获取 MFA qrcode

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /auth/mfa/qrcode:
    get:
      summary: 获取 MFA qrcode
      deprecated: false
      description: 获取绑定 MFA时所需的二维码
      tags:
        - 用户
      parameters:
        - name: Access-Token
          in: header
          description: ''
          required: false
          example: ''
          schema:
            type: string
        - name: Refresh-Token
          in: header
          description: ''
          required: false
          example: ''
          schema:
            type: string
      responses:
        '200':
          description: ''
          content:
            application/json:
              schema:
                type: object
                properties:
                  base: &ref_0
                    $ref: '#/components/schemas/%E5%93%8D%E5%BA%94%E7%8A%B6%E6%80%81'
                  data:
                    type: object
                    properties:
                      secret:
                        type: string
                        title: 密钥
                        description: 多重身份验证密钥
                      qrcode:
                        type: string
                        title: 图片
                        description: 需要 base64 编码
                    x-apifox-orders:
                      - secret
                      - qrcode
                    required:
                      - secret
                      - qrcode
                    x-apifox-ignore-properties: []
                x-apifox-orders:
                  - base
                  - data
                required:
                  - base
                  - data
                x-apifox-refs: {}
                x-apifox-ignore-properties: []
              example:
                base:
                  code: 10000
                  msg: Success
                data:
                  secret: 2dGJT1gjtzo4zNybLa9A
                  qrcode: >-
                    data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAPAAAABQCAMAAAAQlwhOAAAAP1BMVEUAAAB/BXjRV8qrMaTzeeyiKJveZNe9Q7aVG46ABnmJD4L5f/KdI5a8QrWMEoXxd+qLEYTudOegJpmgJpm/RbgEHMs/AAAAAXRSTlMAQObYZgAABPVJREFUeJzsWouu4joMjBFIgBAg+P9/XQFtPH4lfeQUtK2Pri6kiTOTiR2nbNpss802+0kjSkRE3HA4HL6J54+N3paA8eHwK4xfsNr77Cgz4x8h3OGRm6+BTxa4d/wThFGFlpxz/KLGv8KXP7UhrP003TyjgHR/DARTSmpF2Hj5HmGposVBuBwzZnHmne12GpL+qFDQ+iaan69dcnMZm/H1VPChQuKg6P+XaP9ZhRaJyx8+0CnJkOuHmvRST34sbf6Djb3f76mP4Hl8o+Ua4pRYkkReEwRcnW/CfYwr9n60B5ezGM8YLHcgSah1Ncj6gm/8PZ+WiddxOuo5fIkEZS5csl4UJRidenCTpCRKAXjUkZ0TxfM2R8oxh+T1BE5GZK2wSY8jx0Ve0GnAZ53jKIrYaLJwcM+8HA7Y5uOSDmAvwdTDEc8QWG5XFceQacluVxAYcEfIlPgosbNVKoTh2240YYNDVLzacjehcHBCSO21+DzhuBwmuu124xgruJyginwTmRxU4RsoLAJjaPElO7VTWJ9XIgHpXBy5jhTGKpMU6RrgkQytB0SoFE6yZFJHKnsIFHaXZQcNlN9XJKgIimjx8fV6HcP0fr+LaKKkgkvweGYwNn1FZbwXtRx0QuEEjSXG2Pd6HcX4fu8Yg6RJsBJ8n88sbOG0LQDPjNUoMuefXT5ITVLgCYQhgpLatwLuE7DUyxI36WqXKZ8HeqDl2zPGRzMUzv/FF3NKFN7bbRCH1y/nq7P31ciesH4wKYbx+lA4GIi3s/vM6W06FurssNvbmPBgcqHxfi5e2FgwC8dfA9tI6vQmlS+xmxqaBW5wJpFzHXbwh3o4AKOzRdVnsFK6FohQDCcWGJm7Yakff+aORwuwtHIieAivUrKHP3QMtcBD6U6YAgz47Xg8moF1V6QVNowDuGPNzy18GyyGsBrWfT76STbGxtKiGzd9lzAMMyOFcx8MxhUvuuZZ3+BW9oS1DoKBx/6o0eZiwztfIdmEnhz1M9/w7uYdwODAxR7zCq2P15RIbejepz+Ve6byQ2yV+6Ryd6ObDF3HI2CbYCIfS76Px2Nc9DiFsF3CCpzb7SYHp3B9p6ZoYCs8P14WJK7COUE6DssJSw2XCicu5drxTVhEKIFfEvsVUpQwSBEcp689vuDWZvsO9epPxZf9roUJB3VxGNxYPIzi69bgweo2qCnzjBjDUX4qHlcEP8HLF761o8SZLNK3FWF7l7ebJ5YNcyB+xAcFoE7G8Ie0/JcwOmZtDJcuUZJrUhfNIP/I0R5hJ6imcCtizmj0dCWd+K7HL+PIvO6086V8oOl3t/5rvKaEDUxH4SiCeXWgNq0rbPbF287nM+nar7LBJhqp8HOehuPU1iCqvllWUZ7bzy8j+T64lgQmG1Epy5RfC+AvcMPgyczWt3aEiQs4WMFsl8tlIkdNaZrE6m3Y0CLL+5WBCfe+nfxxuTRiLO5sTnsVthv9tflk9xzDfQdnx7UkHOGK44iUwqNde3MlXWqKXgsQHipxm6mqSWAJvsOydLOp8EBruJaDMZQyb3Ms3u+zf3IsgZ1OJwFhiTnFbEbgv538dNKMF91Vy0tsCC9tX1f4v7e18d1ss802+9jnLejX7bnURP177i/b87kU4x8hvDqFF7S18d1sWVtbhcx3oJXEVSa8lkyyOoVXF8ObbbbZ/2H/AgAA//83shRkiJ881wAAAABJRU5ErkJggg==
          headers: {}
          x-apifox-name: 成功
        x-200:失败:
          description: ''
          content:
            application/json:
              schema:
                title: ''
                type: object
                properties:
                  base: *ref_0
                x-apifox-orders:
                  - base
                required:
                  - base
                x-apifox-ignore-properties: []
              example:
                base:
                  code: -1
                  msg: 密码错误
          headers: {}
          x-apifox-name: 失败
      security: []
      x-apifox-folder: 用户
      x-apifox-status: released
      x-run-in-apifox: https://app.apifox.com/web/project/3905938/apis/api-141546915-run
components:
  schemas:
    响应状态:
      type: object
      properties:
        code:
          type: integer
        msg:
          type: string
      x-apifox-orders:
        - code
        - msg
      required:
        - code
        - msg
      x-apifox-ignore-properties: []
      x-apifox-folder: ''
  securitySchemes: {}
servers:
  - url: http://localhost:10001
    description: 开发环境
  - url: localhost:8888
    description: 测试环境
  - url: https://14efdb6874148af54dd6c98f749a9412-app.1024paas.com/douyin
    description: 正式环境
security: []

```
