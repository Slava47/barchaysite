/**
 * Menu data for Li Bo Tea Bar (Ли Бо — Чайный Бар)
 * This file serves as the default/fallback data.
 * In production, data is fetched from the admin API.
 */

const MENU_DATA = {
    categories: [
        { id: 'cold', name: 'Холодные коктейли', nameZh: '冷饮' },
        { id: 'hot', name: 'Горячие коктейли', nameZh: '热饮' },
        { id: 'alco', name: 'Алкогольные коктейли', nameZh: '酒饮' }
    ],

    items: [
        /* ── Холодные коктейли ───────────────────── */
        {
            id: 'c1', category: 'cold',
            name: 'Банановое молоко',
            price: '230',
            description: 'Сладкий, легкий, фруктовый, тягучий, мягкий.',
            fullDescription: 'Сладкий, легкий, фруктовый, тягучий, мягкий.',
            image: 'picture/Холодные коктейли/bananovoe_moloko.jpg',
            images: ['picture/Холодные коктейли/bananovoe_moloko.jpg'],
            tags: ['сладкий', 'фруктовый', 'мягкий', 'лёгкий', 'молочный', 'холодный']
        },
        {
            id: 'c2', category: 'cold',
            name: 'Лиловая гуанинь',
            price: '260',
            description: 'Освежающий, с кислинкой, цветочный, легкий.',
            fullDescription: 'Освежающий, с кислинкой, цветочный, легкий.',
            image: 'picture/Холодные коктейли/lilovaya_guanin.jpg',
            images: ['picture/Холодные коктейли/lilovaya_guanin.jpg'],
            tags: ['освежающий', 'цветочный', 'лёгкий', 'кислый', 'холодный']
        },
        {
            id: 'c3', category: 'cold',
            name: 'Старый князь',
            price: '340',
            description: 'Крепкий, забористый, классический, цитрусовый, яркий.',
            fullDescription: 'Крепкий, забористый, классический, цитрусовый, яркий.',
            image: 'picture/Холодные коктейли/stariy_knyaz.jpg',
            images: ['picture/Холодные коктейли/stariy_knyaz.jpg'],
            tags: ['крепкий', 'классический', 'цитрусовый', 'насыщенный', 'холодный']
        },
        {
            id: 'c4', category: 'cold',
            name: 'Нефритовая река',
            price: '430',
            description: 'Десертный, тягучий, цитрусовый, дождливый.',
            fullDescription: 'Десертный, тягучий, цитрусовый, дождливый.',
            image: 'picture/Холодные коктейли/nefritovaya_reka.jpg',
            images: ['picture/Холодные коктейли/nefritovaya_reka.jpg'],
            tags: ['сладкий', 'десертный', 'цитрусовый', 'мягкий', 'холодный']
        },
        {
            id: 'c5', category: 'cold',
            name: 'Золотая обезьяна',
            price: '340',
            description: 'Сладковато-пряный, десертный, терпкий и фруктовый.',
            fullDescription: 'Сладковато-пряный, десертный, терпкий и фруктовый.',
            image: 'picture/Холодные коктейли/zolotaya_obezyana.jpg',
            images: ['picture/Холодные коктейли/zolotaya_obezyana.jpg'],
            tags: ['пряный', 'десертный', 'терпкий', 'фруктовый', 'холодный']
        },
        {
            id: 'c6', category: 'cold',
            name: 'Тайваньские пираты',
            price: '470',
            description: 'Насыщенный, необычный, ягодный, тягучий, крепкий, темный.',
            fullDescription: 'Насыщенный, необычный, ягодный, тягучий, крепкий, темный.',
            image: 'picture/Холодные коктейли/tayvanskie_piraty.jpg',
            images: ['picture/Холодные коктейли/tayvanskie_piraty.jpg'],
            tags: ['насыщенный', 'ягодный', 'крепкий', 'кислый', 'холодный']
        },
        {
            id: 'c7', category: 'cold',
            name: 'Аметистовое вино',
            price: '320',
            description: 'Космический, свежий, тайский, пикантный, с нежной кислинкой.',
            fullDescription: 'Космический, свежий, тайский, пикантный, с нежной кислинкой.',
            image: 'picture/Холодные коктейли/ametistovoe_vino.jpg',
            images: ['picture/Холодные коктейли/ametistovoe_vino.jpg'],
            tags: ['свежий', 'кислый', 'лёгкий', 'экзотический', 'холодный']
        },
        {
            id: 'c8', category: 'cold',
            name: 'Сестрицы мэй',
            price: '360',
            description: 'Ягодный, яркий, взрывной, с кислинкой, цветочно-садовый.',
            fullDescription: 'Ягодный, яркий, взрывной, с кислинкой, цветочно-садовый.',
            image: 'picture/Холодные коктейли/sestricy_mey.jpg',
            images: ['picture/Холодные коктейли/sestricy_mey.jpg'],
            tags: ['ягодный', 'цветочный', 'кислый', 'яркий', 'холодный']
        },
        {
            id: 'c9', category: 'cold',
            name: 'Южный феникс',
            price: '500',
            description: 'Насыщенный, фруктовый, необычный, бодрящий, ореховый.',
            fullDescription: 'Насыщенный, фруктовый, необычный, бодрящий, ореховый.',
            image: 'picture/Холодные коктейли/yuzhniy_feniks.jpg',
            images: ['picture/Холодные коктейли/yuzhniy_feniks.jpg'],
            tags: ['насыщенный', 'фруктовый', 'бодрящий', 'ореховый', 'холодный']
        },
        {
            id: 'c10', category: 'cold',
            name: 'Цветы и птицы Сюй Вэя',
            price: '320',
            description: 'Освежающий, цветочный, с тонкой сладкой ноткой.',
            fullDescription: 'Освежающий, цветочный, с тонкой сладкой ноткой.',
            image: 'picture/Холодные коктейли/cvety_i_pticy.jpg',
            images: ['picture/Холодные коктейли/cvety_i_pticy.jpg'],
            tags: ['освежающий', 'цветочный', 'сладкий', 'лёгкий', 'холодный']
        },
        {
            id: 'c11', category: 'cold',
            name: 'Гроздья ягод бытия',
            price: '430',
            description: 'Ягодный, текстурный, светлый, с легчайшей терпчинкой.',
            fullDescription: 'Ягодный, текстурный, светлый, с легчайшей терпчинкой.',
            image: 'picture/Холодные коктейли/grozdya_yagod.jpg',
            images: ['picture/Холодные коктейли/grozdya_yagod.jpg'],
            tags: ['ягодный', 'свежий', 'терпкий', 'лёгкий', 'холодный']
        },

        /* ── Горячие коктейли ────────────────────── */
        {
            id: 'h1', category: 'hot',
            name: 'Без тревог',
            price: '350',
            description: 'Чистый, сладостный, фруктовый, успокаивающий.',
            fullDescription: 'Чистый, сладостный, фруктовый, успокаивающий.',
            image: 'picture/Горячие коктейли/bez_trevog.jpg',
            images: ['picture/Горячие коктейли/bez_trevog.jpg'],
            tags: ['сладкий', 'фруктовый', 'мягкий', 'успокаивающий', 'тёплый']
        },
        {
            id: 'h2', category: 'hot',
            name: 'Сычуаньские перцы',
            price: '310',
            description: 'Острый, сытный, коньячный.',
            fullDescription: 'Острый, сытный, коньячный.',
            image: 'picture/Горячие коктейли/sychuanskie_percy.jpg',
            images: ['picture/Горячие коктейли/sychuanskie_percy.jpg'],
            tags: ['пряный', 'острый', 'насыщенный', 'согревающий', 'тёплый']
        },
        {
            id: 'h3', category: 'hot',
            name: 'Красная обезьяна',
            price: '320',
            description: 'Восточный, десертный, сладкий, глубокий.',
            fullDescription: 'Восточный, десертный, сладкий, глубокий.',
            image: 'picture/Горячие коктейли/krasnaya_obezyana.jpg',
            images: ['picture/Горячие коктейли/krasnaya_obezyana.jpg'],
            tags: ['сладкий', 'десертный', 'пряный', 'глубокий', 'тёплый']
        },
        {
            id: 'h4', category: 'hot',
            name: 'Чутка киселе',
            price: '320',
            description: 'Сладкий, фруктовый, согревающий, новогодний.',
            fullDescription: 'Сладкий, фруктовый, согревающий, новогодний.',
            image: 'picture/Горячие коктейли/chutka_kisele.jpg',
            images: ['picture/Горячие коктейли/chutka_kisele.jpg'],
            tags: ['сладкий', 'фруктовый', 'согревающий', 'мягкий', 'тёплый']
        },
        {
            id: 'h5', category: 'hot',
            name: 'Горячая свинюшка',
            price: '490',
            description: 'Копченый, десертный, неожиданный.',
            fullDescription: 'Копченый, десертный, неожиданный.',
            image: 'picture/Горячие коктейли/goryachaya_svinyushka.jpg',
            images: ['picture/Горячие коктейли/goryachaya_svinyushka.jpg'],
            tags: ['дымный', 'десертный', 'сладкий', 'необычный', 'тёплый']
        },
        {
            id: 'h6', category: 'hot',
            name: 'Правила Чэн Ай Сао',
            price: '470',
            description: 'Плотный, тактильный, темный, фруктовый, пиратский.',
            fullDescription: 'Плотный, тактильный, темный, фруктовый, пиратский.',
            image: 'picture/Горячие коктейли/pravila_chen_ay_sao.jpg',
            images: ['picture/Горячие коктейли/pravila_chen_ay_sao.jpg'],
            tags: ['насыщенный', 'фруктовый', 'кислый', 'крепкий', 'тёплый']
        },
        {
            id: 'h7', category: 'hot',
            name: 'Лунный апельсин',
            price: '380',
            description: 'Дымный, терпкий, пряный, цитрусовый.',
            fullDescription: 'Дымный, терпкий, пряный, цитрусовый.',
            image: 'picture/Горячие коктейли/lunniy_apelsin.jpg',
            images: ['picture/Горячие коктейли/lunniy_apelsin.jpg'],
            tags: ['дымный', 'терпкий', 'пряный', 'цитрусовый', 'тёплый']
        },
        {
            id: 'h8', category: 'hot',
            name: 'Еще киселе',
            price: '390',
            description: 'Яркий, кисло-сладкий, на вкус как лето.',
            fullDescription: 'Яркий, кисло-сладкий, на вкус как лето.',
            image: 'picture/Горячие коктейли/eshe_kisele.jpg',
            images: ['picture/Горячие коктейли/eshe_kisele.jpg'],
            tags: ['кислый', 'сладкий', 'фруктовый', 'яркий', 'тёплый']
        },
        {
            id: 'h9', category: 'hot',
            name: 'Осенняя дымка',
            price: '420',
            description: 'Насыщенный, фруктовый, необычный, бодрящий, ореховый.',
            fullDescription: 'Насыщенный, фруктовый, необычный, бодрящий, ореховый.',
            image: 'picture/Горячие коктейли/osennyaya_dymka.jpg',
            images: ['picture/Горячие коктейли/osennyaya_dymka.jpg'],
            tags: ['насыщенный', 'фруктовый', 'бодрящий', 'ореховый', 'тёплый', 'молочный']
        },
        {
            id: 'h10', category: 'hot',
            name: 'Полночь в саду',
            price: '320',
            description: 'Тягучий, плотный, молочный, сладкий, уютный.',
            fullDescription: 'Тягучий, плотный, молочный, сладкий, уютный.',
            image: 'picture/Горячие коктейли/polnoch_v_sadu.jpg',
            images: ['picture/Горячие коктейли/polnoch_v_sadu.jpg'],
            tags: ['сладкий', 'молочный', 'десертный', 'мягкий', 'тёплый']
        },
        {
            id: 'h11', category: 'hot',
            name: 'Чукинский экспресс',
            price: '430',
            description: 'Сладкий, фруктовый, согревающий, новогодний.',
            fullDescription: 'Сладкий, фруктовый, согревающий, новогодний.',
            image: 'picture/Горячие коктейли/chukinskiy_ekspress.jpg',
            images: ['picture/Горячие коктейли/chukinskiy_ekspress.jpg'],
            tags: ['сладкий', 'фруктовый', 'согревающий', 'имбирный', 'тёплый', 'молочный']
        },

        /* ── Алкогольные коктейли ───────────────── */
        {
            id: 'a1', category: 'alco',
            name: 'Биси',
            price: '550',
            description: 'Алкогольный коктейль на основе чая.',
            fullDescription: 'Алкогольный коктейль на основе чая.',
            image: 'picture/Алкогольные коктейли/bisi.jpg',
            images: ['picture/Алкогольные коктейли/bisi.jpg'],
            tags: ['цветочный', 'крепкий', 'травяной', 'алкогольный', 'холодный']
        },
        {
            id: 'a2', category: 'alco',
            name: 'Яцзы',
            price: '550',
            description: 'Алкогольный коктейль на основе чая.',
            fullDescription: 'Алкогольный коктейль на основе чая.',
            image: 'picture/Алкогольные коктейли/yaczi.jpg',
            images: ['picture/Алкогольные коктейли/yaczi.jpg'],
            tags: ['кислый', 'ягодный', 'крепкий', 'алкогольный', 'холодный']
        },
        {
            id: 'a3', category: 'alco',
            name: 'Чивэнь',
            price: '550',
            description: 'Алкогольный коктейль на основе чая.',
            fullDescription: 'Алкогольный коктейль на основе чая.',
            image: 'picture/Алкогольные коктейли/chiven.jpg',
            images: ['picture/Алкогольные коктейли/chiven.jpg'],
            tags: ['кислый', 'свежий', 'лёгкий', 'алкогольный', 'холодный']
        },
        {
            id: 'a4', category: 'alco',
            name: 'Цуню',
            price: '550',
            description: 'Алкогольный коктейль на основе чая.',
            fullDescription: 'Алкогольный коктейль на основе чая.',
            image: 'picture/Алкогольные коктейли/cunyu.jpg',
            images: ['picture/Алкогольные коктейли/cunyu.jpg'],
            tags: ['фруктовый', 'цветочный', 'крепкий', 'алкогольный', 'холодный']
        },
        {
            id: 'a5', category: 'alco',
            name: 'Чаофэн',
            price: '550',
            description: 'Алкогольный коктейль на основе чая.',
            fullDescription: 'Алкогольный коктейль на основе чая.',
            image: 'picture/Алкогольные коктейли/chaofeng.jpg',
            images: ['picture/Алкогольные коктейли/chaofeng.jpg'],
            tags: ['фруктовый', 'освежающий', 'лёгкий', 'алкогольный', 'холодный']
        },
        {
            id: 'a6', category: 'alco',
            name: 'Цзяоту',
            price: '550',
            description: 'Алкогольный коктейль на основе чая.',
            fullDescription: 'Алкогольный коктейль на основе чая.',
            image: 'picture/Алкогольные коктейли/czyaotu.jpg',
            images: ['picture/Алкогольные коктейли/czyaotu.jpg'],
            tags: ['фруктовый', 'сладкий', 'крепкий', 'алкогольный', 'тёплый']
        },
        {
            id: 'a7', category: 'alco',
            name: 'Пулао',
            price: '550',
            description: 'Алкогольный коктейль на основе чая.',
            fullDescription: 'Алкогольный коктейль на основе чая.',
            image: 'picture/Алкогольные коктейли/pulao.jpg',
            images: ['picture/Алкогольные коктейли/pulao.jpg'],
            tags: ['кислый', 'ягодный', 'дымный', 'алкогольный', 'холодный']
        },
        {
            id: 'a8', category: 'alco',
            name: 'Биань',
            price: '550',
            description: 'Алкогольный коктейль на основе чая.',
            fullDescription: 'Алкогольный коктейль на основе чая.',
            image: 'picture/Алкогольные коктейли/bian.jpg',
            images: ['picture/Алкогольные коктейли/bian.jpg'],
            tags: ['цитрусовый', 'имбирный', 'крепкий', 'алкогольный', 'тёплый']
        },
        {
            id: 'a9', category: 'alco',
            name: 'Суаньни',
            price: '550',
            description: 'Алкогольный коктейль на основе чая.',
            fullDescription: 'Алкогольный коктейль на основе чая.',
            image: 'picture/Алкогольные коктейли/suanni.jpg',
            images: ['picture/Алкогольные коктейли/suanni.jpg'],
            tags: ['сладкий', 'десертный', 'молочный', 'алкогольный', 'тёплый']
        }
    ]
};

if (typeof module !== 'undefined' && module.exports) {
    module.exports = MENU_DATA;
}
