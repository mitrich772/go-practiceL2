$(function() {
  
  // Слайдер
  if (document.querySelector('.header-slider')) {
    var splide = new Splide( '.header-slider', {
      type: 'slide',
      perPage: 1,
      arrows: false,
      pagination: true,
      updateOnMove: true,
      gap: 10,
      perMove: 1,
      speed: 800,
    });
    splide.mount();
  }

  // lazyload картинок , берутся из data-src
  const handleLazyImg = () => {
    if ('loading' in HTMLImageElement.prototype) {
      const images = document.querySelectorAll('img[loading="lazy"]');
      images.forEach(img => {
        img.src = img.dataset.src;
      });
    } else {
      const script = document.createElement('script');
      script.src =
        'https://cdnjs.cloudflare.com/ajax/libs/lazysizes/5.1.2/lazysizes.min.js';
      document.body.appendChild(script);
    }
  }
  handleLazyImg();

  // Модальные окна 
  function openModal(){  
    var buttons = document.querySelectorAll('.trigger[data-modal-trigger]');
    var body = document.querySelector('body');

    for(let button of buttons) {
      button.addEventListener('click', () => {
        var trigger = button.getAttribute('data-modal-trigger');
        var modal = document.querySelector(`[data-modal=${trigger}]`);
        var contentWrapper = modal.querySelector('.content-wrapper');
        var close = modal.querySelector('.close');

        close.addEventListener('click', () => {
          modal.classList.remove('open');
          body.classList.remove('blur');
          document.querySelector('.modal-2 .content').innerHTML = '';
        });
        modal.addEventListener('click', () => {
          modal.classList.remove('open');
          body.classList.remove('blur');
          document.querySelector('.modal-2 .content').innerHTML = '';
        });
        contentWrapper.addEventListener('click', (e) => e.stopPropagation());

        modal.classList.toggle('open');
        body.classList.toggle('blur');  
      });
    }
  }
  openModal();

  // получить страницу с товаром и воткнуть в попап  
  const itemBtns = $('.trigger[data-modal-trigger="trigger-2"]');
  itemBtns.each(function(){
    $(this).click(function(e) {
      $('.modal-2 .content').html('<img width="160" height="160" src="/design/img/oval.svg" alt="">');
      $.get( $(this).data('url'), function(data) {
        $('.modal-2 .content').html(data);
        refreshFsLightbox();                  
      }).done(function(data) {
        quantityPlusMinus();
        shopCoupon();
        itemRelated();         
        $(".add-basket.into-product").on('click', function(){
          qty = $(this).parent().find('.input-number-wrapper input[type="number"]').val();
          let href = $(this).attr('href');
          $(this).addClass('active').text('Добавлено').attr('href', 'javascript:void(0)');
          addItemBasket(href, qty);
          return false;
        });
      });
    });
  });  
  
  // подтягиваем товары по каждой категории
  $('.catalog-btns .catalog-content__item').each(function(){
    $(this).appendTo($('.catalog-content'))
  });

  $('.catalog-btns a').click(function(e) {
      e.preventDefault();
      $('.catalog-btns').find('.active').removeClass('active');
      $(this).addClass('active');
      /*$('.catalog-content').find('.catalog-content__item').hide();
      $('#' + $(this).data('switch')).css('display', 'grid');*/
  });

  // Выбрать доп товары + тесто
  function itemRelated(){
    $('.item-related__item').each(function(){
      $(this).on('click', function(){
        if ($(this).hasClass('checked')) {
          $(this).removeClass('checked');
          if ($('.diamtr-choice').length) 
            urlToCart();
        } else {
          $(this).addClass('checked');
          if ($('.diamtr-choice').length) 
            urlToCart();
        }
      });
    });

    // Выбор пицц у трёх
    let step = 3;
    if ($('#pizza17').length) { // если Сытное Комбо
      $('.pizza-item:nth-child(1)').addClass('selected');
      step = 1
    } else if ($('#pizza137').length) {
      $('.pizza-item:nth-child(1)').addClass('selected');
      $('.pizza-item:nth-child(2)').addClass('selected');
      step = 2
    } else {
      $('.pizza-item:nth-child(1), .pizza-item:nth-child(2), .pizza-item:nth-child(3)').addClass('selected');      
    } 
    $('.pizza-item').each(function(){
      $(this).on('click', function(){
          if (!$(this).hasClass('selected') && $('.pizza-item.selected').length < step) {
            $(this).addClass('selected')
          } else {
            $(this).removeClass('selected')
          }
          urlToCart();
      });      
    });   
    
    // Добавим тесто к урлу
    if ($('.testo-choice').length) {      
      $('#diamtrChoice1, #testoChoice12').prop('checked', true);
      $('.tab[for="diamtrChoice1"]').click();
      $('.testo-choice input[name="testo"]').each(function(){
        $(this).on('change', function(){
          urlToCart();
        });
      });
    }
    // Выбор размера
    $('.diamtr-choice input[name="diamtr"]').each(function(){
      $(this).on('change', function(){
        urlToCart();
        $('.item-text__price').text($(this).data('price'));
        let id = $(this).attr('id');
        $('.item-related__price').each(function(){
          $(this).text($(this).data(id));
          $(this).parent().data('url', $(this).data(id+1))
        });
      });
    });
    // Выбор напитка
    if ($('.beverage-choice').length) {
      $('.beverage-choice option:first').attr('selected','selected');
      $('.beverage-choice').on('change', function() {
        urlToCart();
      });
    }
    // Выбор соуса
    if ($('.sauce-choice').length) {
      $('.sauce-choice option:first').attr('selected','selected');
      $('.sauce-choice').on('change', function() {
        urlToCart();
      });
    }
    // Выбор добавок
    if ($('.additives-choice').length) {
      $('.additives-choice option:first').attr('selected','selected');
      $('.additives-choice').on('change', function() {
        urlToCart();
      });
    }
    if ($('.add-basket.into-product').length)
      urlToCart();
  }

  // формируем урл добавления в корзину выбранных комбо  
  function urlToCart(){ 
    const urlBasket = $('.add-basket.into-product').data("href");  
    let beverage = '&0',
        sauce = '&0',
        additives = '&0',
        diamtr = '&0', 
        testo = '&0';
    if ($('.additives-choice').length) {
      additives = '&selected[additives]=' + $('.additives-choice option:selected').val();
    }
    if ($('.sauce-choice').length) {
      sauce = '&selected[sauce]=' + $('.sauce-choice option:selected').val();
    }
    if ($('.beverage-choice').length) {
      beverage = '&selected[beverage]=' + $('.beverage-choice option:selected').val();
    }
    if ($('.diamtr-choice').length) {
      diamtr = '&' + $('.diamtr-choice input[name="diamtr"]:checked').data('url');
    }
    if ($('.testo-choice').length) {
      testo = '&' + $('.testo-choice input[name="testo"]:checked').data('url');
    }
    $('.add-basket.into-product').attr('href', urlBasket + beverage + sauce + additives + testo + diamtr + namePizza() + nameIngredients());
  }

  // выбранные пиццы к урлу
  function namePizza(){
    let pizza = [];
    $('.pizza-item.selected').each(function(){
      let pizzaName = $(this).data('url');
      pizza.push(pizzaName);
    });
    return '&selected[pizza]=' + pizza.join(', ');
  }
  // выбранные ингридиенты к урлу
  function nameIngredients(){
    if ($('.item-related__item').length) {
      let ingredients = [];
      $('.item-related__item.checked').each(function(){
        let ingredientsName = $(this).find('.item-related__name').text();
        ingredients.push(ingredientsName);
      });
      return '&selected[ingredients]=' + ingredients.join(', ');
    } else {
      return '';
    }
  }
  
  // обновление корзины при смене кол-ва товаров
  $(document).on('change keyup', '.cart-items__list input[type="number"]', function(event){
    var key = event.keyCode || event.charCode;
    if( key == 8 || key == 46 ){ 
      return false;
    } else {      
      reloadItems();
      shopCoupon();
      return false;
    }
  }); 

  // обновление корзины при удалении товара
  $(document).on('click', '.cart-items__list .cart-item__remove', function(){
    $(this).parents('.cart-item').hide();    
    $(this).parents('.cart-item__right').find('input[type="number"]').val('0');    
    reloadItems(); 
    shopCoupon(); 
    setTimeout(function() {
      replaceSticker();  
    }, 1000);  
    return false;
  });

  function shopCoupon(){
    // Показать промо-код
    $(document).on('click', '.coupon-header', function(){
      $(this).hide();
      $('.cart-coupon__inner').css('display', 'flex')
    });
    // Обновить купон
    $(document).on('click', '.btn-coupon', function(){    
      reloadItems();
      return false;
    });
    // Удалить купон
    $(document).on('click', '.btn-coupon-delete', function(even){
      $('.cart-items__list input[name="catalog_basket_remove_promocode"]').val('1');
      $.post( $(this).attr('href'), function() {
        reloadItems();
      });      
      return false;
    });
  }
  shopCoupon();

  // Добавляем в корзину товар
  $(".add-basket").on('click', function(){
    let href = $(this).attr('href');
    addItemBasket(href, 1);
    return false;
  });

  function addItemBasket(href, qty){
    $.ajax({
      type: 'POST',
      dataType: 'json',
      context: this,
      url: href + '&qty=' + qty,
      success: function(jsondata){        
        $('.header-cart').addClass('jello-horizontal');
      }, 
      complete: function () {
        setTimeout(function() {
          $('.header-cart').removeClass('jello-horizontal');     
        }, 1000);           
        // тут добавим допы (соусы и прочее)
        if ( $('.item-related__box').length ) {
          $('.item-related__item.checked').each(function() {
            var url = $(this).data("url");
            var arr = $.makeArray( url );  
            for (var i in arr) {
              sendItems(arr[i]);
            }
          });
        }
        setTimeout(function() {
          replaceSticker();
          shopCoupon();          
        }, 1000);
      }
    });      
  }

  // Покажем добавленные товары
  function addedItems(){
    $('.catalog-item .add-basket').each(function(){
      let idItem = $(this).data('id');
      if( $.inArray(idItem, addedItemsInCart()) != -1) {
        $(this).addClass('active').text('Добавлено');
      } else {
        $(this).removeClass('active').text('В корзину');
      }
    });
  }
  addedItems();
  function addedItemsInCart(){ // перебор добавленных товаров
    let items = [];
    $('.cart-items .cart-item__name').each(function() {
      let idItem = $(this).data('id');
      items.push(idItem);
    });
    return items;
  } 
  

  // Обновление корзины аяксом
  let korzina;
  if ($('.order').length) {
    korzina = '/ajax-korziny-zakaza/';
  } else {
    korzina = '/ajax-korziny/';
  }
  function reloadItems() {    
    var msg  = $('.cart-items__list form').serialize();    
    $.ajax({
      type: "post",
      url:  korzina,
      data: msg,
      dataType:"html",
      context: this,
      beforeSend: function (data) {
          $('.cart-items__list').addClass('active');
      },
      success: function (data) {    
        setTimeout(function() {
          $.get( korzina, function( data ) { 
            var fullContent = '<div>' + data + '</div>';
            var html = $(fullContent).find('.cart-items__list form').html();            
            $('.cart-items__list form').html(html);   
            quantityPlusMinus();
            if ($('.cart-item__name.empty').length) {
              $('.cart-items__list .btn').hide()
            }
            if ($('.order').length && $('#suggest').val().length) {
              $('#priceTotal').text(($('#priceTotal').data('price')+zonaPrice).toLocaleString());
              $('#priceDelivery').text(zonaPrice)
            }
          });
        }, 800);        
      }, 
      complete: function (data) {
          setTimeout(function() {
            $('.cart-items__list').removeClass('active');
          }, 800);   
      },
      error: function (xhr, ajaxOptions, thrownError) {
        console.log('Ошибка отправки запроса, получено сообщение:' + xhr.status +' / '+ thrownError);
      }
    });
  }

  // Заменяем аяксом список товаров в стикере корзины
  function replaceSticker() {
    $.ajax({
      url: '/ajax-stiker/',
      type: 'get',
      datatype: 'html',
      success: function(html) {
        $('.header-cart').html(html);
        quantityPlusMinus();
        addedItems();
        posRightCart();
      }
    });
  }

  // функция отправки несколько товаров
  function sendItems(x) { 
      $.ajax({
          url: x
      })
  }

  // плюс-минус в форме заказа товара
  function quantityPlusMinus(){
    const inputWrapperList = document.getElementsByClassName('input-number-wrapper');

    for(let wrapper of inputWrapperList) {
      if (wrapper.querySelector('.increase')) {
        wrapper.querySelector('.increase').remove();
      }
      if (wrapper.querySelector('.decrease')) {
        wrapper.querySelector('.decrease').remove();
      }
      wrapper.insertAdjacentHTML("afterbegin", '<span class="decrease">—</span>');
      wrapper.insertAdjacentHTML("beforeend", '<span class="increase">+</span>');

      const input = wrapper.querySelector('input');
      const incrementation = +input.step || 1;

      wrapper.querySelector('.increase').addEventListener('click', function(even) {
        incInputNumber(input, incrementation);
        event.preventDefault();
      });
      
      wrapper.querySelector('.decrease').addEventListener('click', function(even) {
        incInputNumber(input, "-" + incrementation);
        event.preventDefault();
      });
    }
  } 
  function incInputNumber(input, step) {
    if(!input.disabled) {
      let val = +input.value;

      if (isNaN(val)) val = 1
      val += +step;

      if(input.max && val > input.max) {
        val = input.max;
      } else if (input.min && val < input.min) {
        val = input.min;
      } else if (val < 1) {
        val = 1;
      }

      input.value = val;
      input.setAttribute("value", val);
      reloadItems();
    }
  }  
  quantityPlusMinus();

  // корзину на всю длину смартфона
  function posRightCart(){
    if ($(window).width() < 767) {
      let posCart = $(window).width() - ($('.header-cart').offset().left + $('.header-cart').outerWidth()) - 5;
      $('.cart-items__list').css('right', '-' + posCart + 'px')
    }
  }
  posRightCart();


});

$(window).on('load', function() {
  if ($(window).width() > 991) {
    // Выравниваем высоту в карточках товара
    var teamMax = -1;
    var team_common = $('.catalog-item__name');
    team_common.each(function () {
      var teamHeight = $(this).outerHeight();
      teamMax = teamHeight > teamMax ? teamHeight : teamMax;
    });
    team_common.css('min-height', teamMax);
  }
});

document.addEventListener("DOMContentLoaded", function() {

  // при скроле вниз фиксируем шапку
  let scrollpos = window.scrollY;
  const headerFix = document.querySelector("body");
  let pos = 95;
  if (document.querySelector('.catalog-btns')) {
    pos = document.querySelector('.catalog-btns').offsetTop + 60;
  } 
  window.addEventListener('scroll', function() { 
    scrollpos = window.scrollY;
    /*if (window.matchMedia('screen and (max-width: 576px)').matches)
      pos = 59;*/
    if (scrollpos >= pos) {
      headerFix.classList.add("fixed-top");
    }
    else {
      headerFix.classList.remove("fixed-top");
    }
  });

	// кнопка меню бургер 
	let burgerButton = document.getElementById('nav-icon');
	let mainMenu = document.querySelector('.mainMenu');
  let linkMenu = document.querySelectorAll('.mainMenu > ul > li > a');
  for(let link of linkMenu) {
    link.addEventListener('click', (e) => {
      /*if (window.location.pathname == '/') {
        e.preventDefault();   */   
        mainMenu.classList.remove('act');
        burgerButton.classList.remove('open');
      //}
    });
  }
	burgerButton.addEventListener('click', () => {
		burgerButton.classList.toggle('open');
		if ( burgerButton.classList.contains('open') ) {
			mainMenu.classList.add('act');
		} else {
			mainMenu.classList.remove('act');
		}
	});
	// скрытый инпут против спама формы 
  if ( document.querySelector('.ncapt') ){    
      const input = document.querySelectorAll('.ncapt');
      for(let inp of input) {
        inp.insertAdjacentHTML( 'afterend', '<input type="hidden" name="ncapt" value="' + document.querySelector('.ncapt').textContent + '">' );
      }
  }

  // Отправка форм ajax 
  const ajaxSend = async (formData,formUrl) => {
      const fetchResp = await fetch(formUrl, {
          method: 'POST',
          body: formData
      });
      if (!fetchResp.ok) {
          throw new Error(`Ошибка по адресу ${url}, статус ошибки ${fetchResp.status}`);
      }
      return await fetchResp.text();
  };

  const forms = document.querySelectorAll('form.send-form');

  forms.forEach(form => {    
      form.addEventListener('submit', function (event) {
          event.preventDefault();

          const formData = new FormData(this);
          const formUrl = this.getAttribute('action');
          const length = this.querySelector('input[type="tel"]').value.length;
          
          if (length < 17) {
          	this.querySelector('input[type="tel"]').focus();
          } else {          
          	ajaxSend(formData,formUrl)
              .then((data) => {
                var parser = new DOMParser();
                var doc = parser.parseFromString(data, "text/html");
                  //form.innerHTML = doc.getElementsByClassName("form-send__message")[0].innerHTML;
                  form.innerHTML = '<span class="form-send">...отправляю</span>';
                  setTimeout(function(){
                   		window.location = "/thanks/";
                   	}, 500);
              })
              .catch((err) => console.error(err));
	        }          
          return false;
      });
  });

  

  // скролл к якорю
  document.querySelectorAll("a[href^='#']").forEach(link => {
    link.addEventListener("click", function (e) {
      e.preventDefault();
      let href = this.getAttribute("href").substring(1);
      const scrollTarget = document.getElementById(href);
      //const topOffset = document.querySelector(".header").offsetHeight;
      const topOffset = 85; // если не нужен отступ сверху
      const elementPosition = scrollTarget.getBoundingClientRect().top;
      const offsetPosition = elementPosition - topOffset;

      window.scrollBy({
        top: offsetPosition,
        behavior: "smooth"
      });
      if (this.hash !== "") {
          [].forEach.call(document.querySelectorAll('.mainMenu ul a'), function (el) {
              el.classList.remove('active');
          });
          this.classList.add('active');
      }
    });
  });
	
	// Маска телефона
    [].forEach.call( document.querySelectorAll('input[type="tel"]'), function(input) {
	    var keyCode;
	    function mask(event) {
	        event.keyCode && (keyCode = event.keyCode);
	        var pos = this.selectionStart;
	        /*if (pos < 3) event.preventDefault();*/
	        if (pos < 3) this.setSelectionRange(3, 3);
	        var matrix = "+7 (___) ___ ____",
	            i = 0,
	            def = matrix.replace(/\D/g, ""),
	            val = this.value.replace(/\D/g, ""),
	            new_value = matrix.replace(/[_\d]/g, function(a) {
	                return i < val.length ? val.charAt(i++) || def.charAt(i) : a
	            });
	        i = new_value.indexOf("_");
	        if ( (i == 4 && keyCode == 56) || (i == 4 && keyCode == 104) ) {
	          event.preventDefault();
	        }
	        if (i != -1) {
	            i < 5 && (i = 3);
	            new_value = new_value.slice(0, i)
	        }
	        var reg = matrix.substr(0, this.value.length).replace(/_+/g,
	            function(a) {
	                return "\\d{1," + a.length + "}"
	            }).replace(/[+()]/g, "\\$&");
	        reg = new RegExp("^" + reg + "$");
	        if (!reg.test(this.value) || this.value.length < 5 || keyCode > 47 && keyCode < 58) this.value = new_value;
	        if (event.type == "blur" && this.value.length < 5)  this.value = ""
	    }
	    input.addEventListener("input", mask, false);
	    input.addEventListener("focus", mask, false);
	    input.addEventListener("blur", mask, false);
	    input.addEventListener("keydown", mask, false)
	  });

  // проверка браузера на поддержку webp
  function testWebP(elem) {
    const webP = new Image();
    webP.src = 'data:image/webp;base64,UklGRjoAAABXRUJQVlA4IC4AAACyAgCdASoCAAIALmk0mk0iIiIiIgBoSygABc6WWgAA/veff/0PP8bA//LwYAAA';
    webP.onload = webP.onerror = function () {
      webP.height === 2 ? elem.classList.add('webp-true') : elem.classList.add('webp-false')
    }
  }
  testWebP(document.body);
  setTimeout(() => {
    if (document.body.classList.contains('webp-false')) {
      document.querySelectorAll('img[src$=".webp"]').forEach(item => {
        item.src = item.src.replace('.webp', '.jpg');
      });
    }
  }, 500);

  // загрузим карту позже по скролу
  const handleMap = () => {
      const map = document.getElementById('yamap')

      if (map) {
          let ok = false
          window.addEventListener('scroll', function () {
              if (ok === false) {
                  ok = true
                  setTimeout(() => {
                      let script = document.createElement('script')
                      script.src =
                          'https://api-maps.yandex.ru/services/constructor/1.0/js/?um=constructor%3A04f099440f0a3c06f3f638f774f926534c0313abb77e9bc13fb7a98c2c165d04&amp;width=100%25&amp;height=400&amp;lang=ru_RU&amp;scroll=false'
                      document.getElementById('yamap').replaceWith(script)
                  }, 500)
              }
          })
          // блок с контактами на карте выравниваем под .container, т.к. сама карта width=100%
          if (window.matchMedia('screen and (min-width: 768px)').matches) {
              let leftOffset = document.querySelector(".container").getBoundingClientRect().left;
              document.querySelector(".contacts-inner").style.left = leftOffset + 15 + 'px';
          }
      }
  }
  handleMap();

});
